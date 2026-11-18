package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// aptMu serializes apt-get invocations. Without this, two concurrent
// php.ext.apply calls collide on /var/lib/dpkg/lock-frontend and the second
// blocks up to the apt timeout. Held only around the apt subprocess; phpenmod
// and systemctl run unserialized.
var aptMu sync.Mutex

// minimalEnv is the env passed to every subprocess — we do NOT inherit the
// parent env so a compromised caller can't smuggle LD_PRELOAD/PATH tricks
// through apt or systemctl. DEBIAN_FRONTEND stops apt from opening a dialog
// on package upgrades.
var minimalEnv = []string{
	"PATH=/usr/sbin:/usr/bin:/sbin:/bin",
	"DEBIAN_FRONTEND=noninteractive",
	"LC_ALL=C",
}

// Subprocess wrappers exposed as package-level vars so tests inject fakes.
// Every default* runs a real binary; tests swap these out with a defer-restore
// pattern and NEVER touch the host.

var runAptGet = defaultRunAptGet
var runPhpenmod = defaultRunPhpenmod
var runPhpdismod = defaultRunPhpdismod
var runSystemctl = defaultRunSystemctl
var runDpkgQuery = defaultRunDpkgQuery
var globConfD = defaultGlobConfD

// listInstalledPHPVersionsFunc wraps the package's listInstalledPHPVersions
// so apply/list can validate "version is installed" under a fake in tests
// without editing php_version_list.go.
var listInstalledPHPVersionsFunc = listInstalledPHPVersions

// defaultRunAptGet runs `apt-get -y <action> <pkgs...>`. Caller holds aptMu.
// gosec G204: action is a const + pkgs come from phpext.ResolvePackages which
// only yields validated allowlist entries. No user-controlled strings reach argv.
func defaultRunAptGet(ctx context.Context, action string, pkgs ...string) ([]byte, error) {
	if osinfo.Detect().IsRHEL() {
		args := append([]string{"-y", action}, pkgs...)
		cmd := execCommandContext(ctx, "dnf", args...) //nolint:gosec // validated upstream
		cmd.Env = minimalEnv
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return out.Bytes(), err
	}
	args := append([]string{"-y", action}, pkgs...)
	cmd := execCommandContext(ctx, "apt-get", args...) //nolint:gosec // validated upstream
	cmd.Env = minimalEnv
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func defaultRunPhpenmod(ctx context.Context, version, module string) ([]byte, error) {
	if info := osinfo.Detect(); info.IsRHEL() {
		return toggleRemiPHPIni(info.PHPConfDDir(version, "fpm"), module, true)
	}
	cmd := execCommandContext(ctx, "phpenmod", "-v", version, module) //nolint:gosec // validated upstream
	cmd.Env = minimalEnv
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func defaultRunPhpdismod(ctx context.Context, version, module string) ([]byte, error) {
	if info := osinfo.Detect(); info.IsRHEL() {
		return toggleRemiPHPIni(info.PHPConfDDir(version, "fpm"), module, false)
	}
	cmd := execCommandContext(ctx, "phpdismod", "-v", version, module) //nolint:gosec // validated upstream
	cmd.Env = minimalEnv
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func toggleRemiPHPIni(confDir, module string, enable bool) ([]byte, error) {
	if enable {
		disabled, _ := filepath.Glob(filepath.Join(confDir, "*-"+module+".ini.disabled"))
		disabledPlain, _ := filepath.Glob(filepath.Join(confDir, module+".ini.disabled"))
		for _, p := range append(disabled, disabledPlain...) {
			target := strings.TrimSuffix(p, ".disabled")
			if err := os.Rename(p, target); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	enabled, _ := filepath.Glob(filepath.Join(confDir, "*-"+module+".ini"))
	enabledPlain, _ := filepath.Glob(filepath.Join(confDir, module+".ini"))
	for _, p := range append(enabled, enabledPlain...) {
		if err := os.Rename(p, p+".disabled"); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// defaultRunSystemctl runs `systemctl <args...>`. Used for reload + list-units.
// gosec G204: args come from hardcoded strings in reloadFPMs + a version that
// was validated via phpext.ValidVersion + unit names derived from systemctl's
// own output. No user-controlled data in argv.
func defaultRunSystemctl(ctx context.Context, args ...string) ([]byte, error) {
	cmd := execCommandContext(ctx, "systemctl", args...) //nolint:gosec // validated upstream
	cmd.Env = minimalEnv
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// defaultRunDpkgQuery runs `dpkg-query -W -f='${Package}\t${Status}\n' <pattern>`
// and returns the raw text. Callers parse it themselves.
// gosec G204: pattern is always `php<v>-*` where v was validated via phpext.ValidVersion.
func defaultRunDpkgQuery(ctx context.Context, pattern string) ([]byte, error) {
	if osinfo.Detect().IsRHEL() {
		cmd := execCommandContext(ctx, "rpm", "-qa", "--qf", "%{NAME}\tinstall ok installed\n", pattern) //nolint:gosec // validated upstream
		cmd.Env = minimalEnv
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil && out.Len() == 0 && stderr.Len() > 0 {
			return nil, fmt.Errorf("rpm -qa: %w: %s", err, stderr.String())
		}
		return out.Bytes(), nil
	}
	cmd := execCommandContext(ctx, "dpkg-query", "-W", "-f=${Package}\t${Status}\n", pattern) //nolint:gosec // validated upstream
	cmd.Env = minimalEnv
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	// dpkg-query exits non-zero when no packages match the pattern — not a
	// real failure for our purposes. Only surface an error if stdout is also
	// empty AND stderr signals something the caller can't infer from "no rows".
	if err != nil && out.Len() == 0 && !strings.Contains(stderr.String(), "no packages found") {
		return nil, fmt.Errorf("dpkg-query: %w: %s", err, stderr.String())
	}
	return out.Bytes(), nil
}

// defaultGlobConfD lists enabled module ini symlinks for (version, sapi). Each
// entry ends with `<module>.ini`. Returns the matched paths.
func defaultGlobConfD(version, sapi string) ([]string, error) {
	return filepath.Glob(filepath.Join(osinfo.Detect().PHPConfDDir(version, sapi), "*.ini"))
}

// truncateErrorOutput clips subprocess output to at most 512 bytes so a verbose
// apt failure doesn't overflow the wire response. Truncation keeps the TAIL —
// apt/dpkg errors appear at the end of output (post-download, during unpack or
// trigger processing), so the head is success noise and the tail is the verdict.
func truncateErrorOutput(b []byte) string {
	const max = 512
	s := strings.TrimSpace(string(b))
	if len(s) <= max {
		return s
	}
	return "…" + s[len(s)-max:]
}

// hasHardAptError returns true if output contains a line starting with "E: " —
// apt's convention for fatal errors (e.g. "E: Unable to correct problems, you
// have held broken packages"). Trigger warnings and debconf notes don't match.
// Used to flag operator-visible problems in the last_error response field even
// when the verdict readback otherwise matches intent.
func hasHardAptError(output []byte) bool {
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "E: ") {
			return true
		}
	}
	return false
}

// reloadFPMs asks systemd to reload php<v>-fpm.service (if present) plus every
// running dotnext-fpm@<user>.service whose user has <v> pinned in the per-user
// version file. Individual reload failures are logged but never block the
// overall apply response — the apt/phpen* mutation already happened on disk,
// and the next reconciler pass or manual reload recovers.
var reloadFPMs = defaultReloadFPMs

func defaultReloadFPMs(ctx context.Context, version string) {
	// Global php<v>-fpm.service. Ignore error — may be masked post-ADR-0025.
	_, _ = runSystemctl(ctx, "reload", osinfo.Detect().PHPFPMServiceUnit(version))

	// Per-user dotnext-fpm instances. list-units output:
	//   dotnext-fpm@alice.service loaded active running ...
	out, err := runSystemctl(ctx, "list-units", "dotnext-fpm@*.service", "--state=running", "--no-legend", "--plain", "--no-pager")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		unit := fields[0]
		// Extract user: dotnext-fpm@alice.service → alice
		at := strings.Index(unit, "@")
		dot := strings.LastIndex(unit, ".service")
		if at < 0 || dot <= at+1 {
			continue
		}
		user := unit[at+1 : dot]
		if !userPinnedToVersion(user, version) {
			continue
		}
		_, _ = runSystemctl(ctx, "reload", unit)
	}
}

// userPinnedToVersion reads /etc/dotnext/user-phpver/<user> and returns
// whether its content (stripped) equals version.
var userPinnedToVersion = func(user, version string) bool {
	p := fmt.Sprintf("/etc/dotnext/user-phpver/%s", user)
	b, err := readSmallFile(p)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(b)) == version
}

// readSmallFile is injectable for tests.
var readSmallFile = os.ReadFile
