package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// systemAptCheckResponse is the wire shape for system.apt_check.
//
// JAB-10: on a healthy host Error is nil and Packages/Total are populated.
// When apt fails (dpkg lock held by a running upgrade, repo unreachable,
// bad perms, unparseable output, …) the handler does NOT return a bare
// agent error — it returns this response with a *structured* Error the
// admin UI can render (reason + failed command + exit code + stderr
// summary + a recovery hint). A hard agent error is reserved for the
// agent being unreachable, which is a different failure than "apt said no".
type systemAptCheckResponse struct {
	Packages       []aptUpgradablePackage `json:"packages"`
	Total          int                    `json:"total"`
	SecurityTotal  int                    `json:"security_total"`
	InstalledTotal int                    `json:"installed_total"`
	Error          *aptCheckError         `json:"error,omitempty"`
	// JAB-353 OS-patch status (cheap file stats, independent of the apt-get
	// update above): RebootRequired mirrors /var/run/reboot-required;
	// LastAppliedAt is when unattended-upgrades last completed an apply.
	RebootRequired     bool       `json:"reboot_required"`
	RebootRequiredPkgs []string   `json:"reboot_required_pkgs,omitempty"`
	LastAppliedAt      *time.Time `json:"last_applied_at,omitempty"`
}

// aptCheckError is the admin-facing, structured diagnostic for a failed
// system-package check (JAB-10).
type aptCheckError struct {
	// Reason is a stable machine key the UI switches on:
	//   apt_locked | repo_unreachable | permission | command_failed
	Reason string `json:"reason"`
	// Command is the failed command, e.g. "apt-get update".
	Command string `json:"command"`
	// ExitCode is the process exit code, or -1 when it never ran / was killed.
	ExitCode int `json:"exit_code"`
	// Stderr is a trimmed, size-capped stderr/stdout summary.
	Stderr string `json:"stderr"`
	// Hint is the next recovery step, in plain admin language.
	Hint string `json:"hint"`
}

type aptUpgradablePackage struct {
	Name    string `json:"name"`
	Current string `json:"current"`
	New     string `json:"new"`
	Source  string `json:"source"`
	// Security is true when the upgrade candidate comes from a *-security
	// suite (e.g. "stable-security"). Coarse severity, zero extra apt calls.
	Security bool `json:"security"`
}

// systemAptCheckHandler runs `apt-get update` then `apt list --upgradable`
// and parses the column-stable output. Every apt invocation includes
// `-o DPkg::Lock::Timeout=60` because unattended-upgrades.timer (and the
// panel's own dotnext-apt-oneshot.service) can hold the dpkg lock; without
// the timeout, ops see a cryptic crash. LC_ALL=C pins the column locale.
//
// This command takes NO user-controlled parameters — apt is invoked with
// fixed args only. Any future params should be type-tagged enums, not
// passed through as strings.
func systemAptCheckHandler(ctx context.Context, _ json.RawMessage) (any, error) {
	if osinfo.Detect().IsRHEL() {
		return systemDnfCheckHandler(ctx)
	}
	if _, stderr, code, err := runAptCapture(ctx, "update", "-qq"); err != nil {
		return systemAptCheckResponse{
			Error: classifyAptError("apt-get update", stderr, code, err),
		}, nil
	}
	out, stderr, code, err := aptListCapture(ctx)
	if err != nil {
		return systemAptCheckResponse{
			Error: classifyAptError("apt list --upgradable", stderr, code, err),
		}, nil
	}
	pkgs := parseAptUpgradable(string(out))
	sec := 0
	for _, p := range pkgs {
		if p.Security {
			sec++
		}
	}
	rebootReq, rebootPkgs := aptRebootRequired()
	return systemAptCheckResponse{
		Packages:           pkgs,
		Total:              len(pkgs),
		SecurityTotal:      sec,
		InstalledTotal:     countInstalledPackages(ctx),
		RebootRequired:     rebootReq,
		RebootRequiredPkgs: rebootPkgs,
		LastAppliedAt:      aptLastApplied(),
	}, nil
}

func systemDnfCheckHandler(ctx context.Context) (systemAptCheckResponse, error) {
	out, stderr, code, err := dnfCheckUpdateCapture(ctx)
	if err != nil {
		return systemAptCheckResponse{
			Error: classifyDnfError("dnf check-update", stderr, code, err),
		}, nil
	}
	secSet := dnfSecurityUpdates(ctx)
	currentVersions := rpmInstalledVersions(ctx)
	pkgs := parseDnfCheckUpdate(string(out), currentVersions, secSet)
	sec := 0
	for _, p := range pkgs {
		if p.Security {
			sec++
		}
	}
	rebootReq, rebootPkgs := dnfRebootRequired(ctx)
	return systemAptCheckResponse{
		Packages:           pkgs,
		Total:              len(pkgs),
		SecurityTotal:      sec,
		InstalledTotal:     countInstalledRPMs(ctx),
		RebootRequired:     rebootReq,
		RebootRequiredPkgs: rebootPkgs,
		LastAppliedAt:      dnfLastApplied(),
	}, nil
}

func dnfCheckUpdateCapture(ctx context.Context) (stdout []byte, stderr string, exitCode int, err error) {
	cmd := execCommandContext(ctx, "dnf", "check-update", "-q")
	cmd.Env = append(cmd.Environ(), "LC_ALL=C")
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := cmd.Run()
	code := exitCodeOf(runErr)
	// dnf check-update exits 100 when updates exist and 0 when up-to-date.
	if runErr != nil && code != 100 {
		combined := strings.TrimSpace(errBuf.String() + "\n" + outBuf.String())
		return nil, combined, code, runErr
	}
	return []byte(outBuf.String()), "", 0, nil
}

func dnfSecurityUpdates(ctx context.Context) map[string]bool {
	out, err := execCommandContext(ctx, "dnf", "updateinfo", "list", "updates", "security", "-q").Output()
	if err != nil {
		return nil
	}
	return parseDnfSecurityUpdates(string(out))
}

func parseDnfSecurityUpdates(out string) map[string]bool {
	res := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// Example row: RHSA-2025:1234 Moderate/Sec. curl-7.76.1-29.el9.x86_64
		nevra := fields[len(fields)-1]
		if dot := strings.LastIndexByte(nevra, '.'); dot > 0 {
			nevra = nevra[:dot]
		}
		// Trim -version-release from N-V-R.
		if dash := strings.LastIndexByte(nevra, '-'); dash > 0 {
			nv := nevra[:dash]
			if dash2 := strings.LastIndexByte(nv, '-'); dash2 > 0 {
				res[nv[:dash2]] = true
			}
		}
	}
	return res
}

func rpmInstalledVersions(ctx context.Context) map[string]string {
	out, err := execCommandContext(ctx, "rpm", "-qa", "--qf", "%{NAME}\t%{VERSION}-%{RELEASE}\n").Output()
	if err != nil {
		return nil
	}
	res := map[string]string{}
	for _, line := range strings.Split(string(out), "\n") {
		name, ver, ok := strings.Cut(line, "\t")
		if ok && name != "" {
			res[name] = strings.TrimSpace(ver)
		}
	}
	return res
}

func parseDnfCheckUpdate(out string, current map[string]string, secSet map[string]bool) []aptUpgradablePackage {
	var pkgs []aptUpgradablePackage
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "Last metadata") || strings.HasPrefix(line, "Obsoleting") || strings.HasPrefix(line, "Security:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		nameArch, newVer, repo := fields[0], fields[1], fields[2]
		name := nameArch
		if dot := strings.LastIndexByte(nameArch, '.'); dot > 0 {
			name = nameArch[:dot]
		}
		if name == "" || newVer == "" {
			continue
		}
		pkgs = append(pkgs, aptUpgradablePackage{
			Name:     name,
			Current:  current[name],
			New:      newVer,
			Source:   repo,
			Security: secSet[name],
		})
	}
	return pkgs
}

func countInstalledRPMs(ctx context.Context) int {
	out, err := execCommandContext(ctx, "rpm", "-qa", "--qf", "1\n").Output()
	if err != nil {
		return 0
	}
	return strings.Count(string(out), "\n")
}

func dnfRebootRequired(ctx context.Context) (bool, []string) {
	cmd := execCommandContext(ctx, "needs-restarting", "-r")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return false, nil
	}
	if exitCodeOf(err) != 1 {
		return false, nil
	}
	var pkgs []string
	seen := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "* "))
		if line == "" || strings.HasPrefix(line, "Core libraries") || strings.HasPrefix(line, "Reboot is required") || strings.HasPrefix(line, "More information") {
			continue
		}
		if !seen[line] {
			seen[line] = true
			pkgs = append(pkgs, line)
		}
	}
	return true, pkgs
}

func dnfLastApplied() *time.Time {
	out, err := execCommand("rpm", "-qa", "--last").Output()
	if err != nil {
		return nil
	}
	line, _, _ := strings.Cut(string(out), "\n")
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	// Strip the leading package name token and parse RPM --last timestamp.
	ts := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
	if i, err := strconv.ParseInt(ts, 10, 64); err == nil && i > 0 {
		t := time.Unix(i, 0).UTC()
		return &t
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, "Mon 02 Jan 2006 03:04:05 PM MST", "Mon Jan 2 15:04:05 2006"} {
		if t, err := time.Parse(layout, ts); err == nil {
			u := t.UTC()
			return &u
		}
	}
	return nil
}

var (
	dnfLockPatterns = []string{
		"waiting for process with pid", "another app is currently holding",
		"could not lock", "cache lock",
	}
	dnfRepoPatterns = []string{
		"failed to download metadata", "cannot prepare internal mirrorlist",
		"curl error", "could not resolve host", "connection timed out",
		"no more mirrors to try",
	}
)

func classifyDnfError(command, stderr string, exitCode int, err error) *aptCheckError {
	low := strings.ToLower(stderr)
	summary := summarizeStderr(stderr)
	if summary == "" && err != nil {
		summary = err.Error()
	}
	switch {
	case containsAny(low, dnfLockPatterns):
		return &aptCheckError{
			Reason: "apt_locked", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "A package operation is in progress (dnf-automatic or a dotNext system update holds the dnf lock). Wait for it to finish, then retry.",
		}
	case containsAny(low, dnfRepoPatterns):
		return &aptCheckError{
			Reason: "repo_unreachable", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "dnf could not reach a repository. Check the server's network/DNS and /etc/yum.repos.d/, then retry.",
		}
	case containsAny(low, aptPermPatterns):
		return &aptCheckError{
			Reason: "permission", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "The package check must run with root privileges. Verify the panel-agent is running as root.",
		}
	default:
		return &aptCheckError{
			Reason: "command_failed", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "dnf returned an error. Run `sudo dnf check-update` on the host to inspect the full output.",
		}
	}
}

// aptRebootRequired reports whether the host needs a reboot to finish applying
// updates (Debian/Ubuntu drop /var/run/reboot-required after a kernel/libc/…
// upgrade) and, best-effort, the packages that triggered it.
func aptRebootRequired() (bool, []string) {
	if _, err := os.Stat("/var/run/reboot-required"); err != nil {
		return false, nil
	}
	var pkgs []string
	if b, err := os.ReadFile("/var/run/reboot-required.pkgs"); err == nil {
		seen := map[string]bool{}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !seen[line] {
				seen[line] = true
				pkgs = append(pkgs, line)
			}
		}
	}
	return true, pkgs
}

// aptLastApplied returns when unattended-upgrades last completed an apply, taken
// from the mtime of the periodic stamp it rewrites after each run. Nil when it
// has never run (fresh host, or the feature was only just enabled).
func aptLastApplied() *time.Time {
	for _, p := range []string{
		"/var/lib/apt/periodic/unattended-upgrades-stamp",
		"/var/lib/apt/periodic/upgrade-stamp",
	} {
		if fi, err := os.Stat(p); err == nil {
			t := fi.ModTime().UTC()
			return &t
		}
	}
	return nil
}

// aptLockPatterns / aptRepoPatterns / aptPermPatterns are matched
// case-insensitively against the failed command's output to pick a Reason.
// Order matters: lock is checked before the generic fallback so a lock held
// by a running upgrade reads as "in progress", not "command failed".
var (
	aptLockPatterns = []string{
		"could not get lock", "unable to acquire the dpkg", "dpkg frontend lock",
		"dpkg::lock::timeout", "lock::timeout", "dpkg was interrupted",
		"is another process using it", "resource temporarily unavailable",
	}
	aptRepoPatterns = []string{
		"failed to fetch", "could not resolve", "temporary failure resolving",
		"could not connect", "unable to connect", "connection failed",
		"connection timed out", "network is unreachable", "cannot initiate the connection",
	}
	aptPermPatterns = []string{
		"permission denied", "are you root", "requires superuser",
		"must be run as root", "eacces", "operation not permitted",
	}
)

// classifyAptError maps a failed apt invocation to a structured, admin-facing
// diagnostic. It never returns nil.
func classifyAptError(command, stderr string, exitCode int, err error) *aptCheckError {
	low := strings.ToLower(stderr)
	summary := summarizeStderr(stderr)
	if summary == "" && err != nil {
		summary = err.Error()
	}
	switch {
	case containsAny(low, aptLockPatterns):
		return &aptCheckError{
			Reason: "apt_locked", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "A package operation is in progress (unattended-upgrades or a DotNext system update holds the dpkg lock). Wait for it to finish, then retry.",
		}
	case containsAny(low, aptRepoPatterns):
		return &aptCheckError{
			Reason: "repo_unreachable", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "apt could not reach a repository. Check the server's network/DNS and its apt sources, then retry.",
		}
	case containsAny(low, aptPermPatterns):
		return &aptCheckError{
			Reason: "permission", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "The package check must run with root privileges. Verify the panel-agent is running as root.",
		}
	default:
		return &aptCheckError{
			Reason: "command_failed", Command: command, ExitCode: exitCode, Stderr: summary,
			Hint: "The package manager returned an error. Run `sudo apt-get update` on the host to reproduce and inspect the full output.",
		}
	}
}

func containsAny(haystack string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

// summarizeStderr trims apt's output to a compact, admin-safe summary: the
// last few non-empty lines (apt puts the actionable error last), capped in
// length so it stays displayable and never leaks a wall of text.
func summarizeStderr(s string) string {
	const maxLines, maxLen = 4, 600
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(l); t != "" {
			lines = append(lines, t)
		}
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	out := strings.Join(lines, "\n")
	if len(out) > maxLen {
		out = out[:maxLen] + "…"
	}
	return out
}

// exitCodeOf extracts the process exit code from an exec error, or -1 when
// the command could not run (binary missing, context cancelled, killed).
func exitCodeOf(err error) int {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

// countInstalledPackages returns the number of installed dpkg packages, for
// the "X of N packages" stat. Best-effort: 0 on error (the UI just omits it).
func countInstalledPackages(ctx context.Context) int {
	out, err := execCommandContext(ctx, "dpkg-query", "-f", "${db:Status-Abbrev}\n", "-W").Output()
	if err != nil {
		return 0
	}
	n := 0
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "ii") {
			n++
		}
	}
	return n
}

// runAptCapture runs apt-get with the shared lock-timeout + locale env and
// returns stdout, a combined output summary source, the exit code, and err.
func runAptCapture(ctx context.Context, args ...string) (stdout []byte, stderr string, exitCode int, err error) {
	full := append([]string{"-o", "DPkg::Lock::Timeout=60"}, args...)
	cmd := execCommandContext(ctx, "apt-get", full...)
	cmd.Env = append(cmd.Environ(), "LC_ALL=C", "DEBIAN_FRONTEND=noninteractive")
	// apt-get writes most diagnostics to stderr but some to stdout; combine
	// so classification sees everything regardless of stream.
	out, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return nil, string(out), exitCodeOf(runErr), runErr
	}
	return out, "", 0, nil
}

// aptListCapture runs `apt list --upgradable`, capturing stderr separately so
// a failure (rare — lock/perm) is classifiable.
func aptListCapture(ctx context.Context) (stdout []byte, stderr string, exitCode int, err error) {
	cmd := execCommandContext(ctx, "apt", "list", "--upgradable")
	cmd.Env = append(cmd.Environ(), "LC_ALL=C")
	var errBuf strings.Builder
	cmd.Stderr = &errBuf
	out, runErr := cmd.Output()
	if runErr != nil {
		return nil, errBuf.String(), exitCodeOf(runErr), runErr
	}
	return out, "", 0, nil
}

// parseAptUpgradable reads `apt list --upgradable` output, format example:
//
//	Listing...
//	curl/stable 8.5.0-2 amd64 [upgradable from: 8.4.0-2]
//	libc6/stable 2.36-9+deb12u4 amd64 [upgradable from: 2.36-9+deb12u3]
var aptListLine = regexp.MustCompile(`^([A-Za-z0-9.+\-]+)/(\S+)\s+(\S+)\s+(\S+)\s+\[upgradable from:\s*(\S+?)\]$`)

func parseAptUpgradable(out string) []aptUpgradablePackage {
	var pkgs []aptUpgradablePackage
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "Listing") || strings.HasPrefix(line, "WARNING:") {
			continue
		}
		m := aptListLine.FindStringSubmatch(line)
		if len(m) != 6 {
			continue
		}
		pkgs = append(pkgs, aptUpgradablePackage{
			Name:     m[1],
			Source:   m[2],
			New:      m[3],
			Current:  m[5],
			Security: strings.Contains(strings.ToLower(m[2]), "security"),
		})
	}
	return pkgs
}

func init() {
	Default.Register("system.apt_check", systemAptCheckHandler)
}
