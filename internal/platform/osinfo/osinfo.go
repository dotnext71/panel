// Package osinfo provides single-source-of-truth host operating system
// detection and distro-specific layout mappings across Debian/Ubuntu and
// Enterprise Linux (AlmaLinux 8, 9, and 10).
package osinfo

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Family identifies the broad Linux distribution family.
type Family string

const (
	FamilyDebian  Family = "debian"
	FamilyRHEL    Family = "rhel"
	FamilyUnknown Family = "unknown"
)

// Info describes the parsed /etc/os-release identity and derived family metadata.
type Info struct {
	ID           string
	IDLike       []string
	VersionID    string
	MajorVersion int
	PrettyName   string
	Codename     string
	Family       Family
}

var osReleasePath = "/etc/os-release"

// Detect reads osReleasePath (default /etc/os-release) and returns the parsed Info.
// If the file is missing or unreadable (e.g., on non-Linux dev hosts), it returns
// FamilyDebian defaults so cross-OS unit tests without fixtures keep deterministic
// behavior unless explicitly overridden.
func Detect() Info {
	return DetectFrom(osReleasePath)
}

// DetectFrom reads path and parses its os-release fields.
func DetectFrom(path string) Info {
	data, err := os.ReadFile(path)
	if err != nil {
		return Info{Family: FamilyDebian}
	}
	return ParseOSRelease(string(data))
}

// ParseOSRelease parses key=value pairs in os-release format.
func ParseOSRelease(content string) Info {
	var info Info
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		val := strings.TrimSpace(v)
		val = strings.Trim(val, `"'`)
		switch k {
		case "ID":
			info.ID = strings.ToLower(val)
		case "ID_LIKE":
			for _, tok := range strings.Fields(strings.ToLower(val)) {
				info.IDLike = append(info.IDLike, tok)
			}
		case "VERSION_ID":
			info.VersionID = val
		case "PRETTY_NAME":
			info.PrettyName = val
		case "VERSION_CODENAME":
			info.Codename = strings.ToLower(val)
		}
	}
	info.MajorVersion = parseMajorVersion(info.VersionID)
	info.Family = classifyFamily(info.ID, info.IDLike)
	return info
}

func parseMajorVersion(versionID string) int {
	if versionID == "" {
		return 0
	}
	majorStr, _, _ := strings.Cut(versionID, ".")
	n, err := strconv.Atoi(strings.TrimSpace(majorStr))
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func classifyFamily(id string, idLike []string) Family {
	switch id {
	case "debian", "ubuntu", "linuxmint", "pop":
		return FamilyDebian
	case "almalinux", "rocky", "rhel", "centos", "ol", "fedora", "virtuozzo", "cloudlinux":
		return FamilyRHEL
	}
	for _, like := range idLike {
		switch like {
		case "debian", "ubuntu":
			return FamilyDebian
		case "rhel", "centos", "fedora":
			return FamilyRHEL
		}
	}
	return FamilyUnknown
}

// IsRHEL reports whether the OS belongs to the RHEL / AlmaLinux family.
func (i Info) IsRHEL() bool {
	return i.Family == FamilyRHEL
}

// IsDebian reports whether the OS belongs to the Debian / Ubuntu family.
func (i Info) IsDebian() bool {
	return i.Family == FamilyDebian
}

// Ubuntu releases this project supports or has verified. Ubuntu version IDs
// are "YY.MM" (22.04, 24.04, 26.04), so MajorVersion is the two-digit year and
// VersionID must be compared whole — "26" alone would also match 26.10.
const (
	UbuntuJammy    = "22.04"
	UbuntuNoble    = "24.04"
	UbuntuResolute = "26.04"
)

// IsUbuntu reports whether this host runs Ubuntu specifically (as opposed to
// a Debian derivative that shares the family).
func (i Info) IsUbuntu() bool {
	return i.ID == "ubuntu"
}

// IsUbuntuRelease reports whether this host runs Ubuntu at exactly the given
// VERSION_ID (e.g. UbuntuNoble).
func (i Info) IsUbuntuRelease(versionID string) bool {
	return i.IsUbuntu() && i.VersionID == versionID
}

// IsSupportedUbuntu reports true for the Ubuntu LTS releases this project
// explicitly supports: 22.04, 24.04 and 26.04. Other Ubuntu versions still
// install (they are Debian-family) but are untested — interim releases are
// deliberately excluded: 25.10 reached end of life in July 2026 and neither
// the ondrej PPA nor packages.sury.org publishes a suite for it, so the PHP
// source would not resolve at all.
func (i Info) IsSupportedUbuntu() bool {
	if !i.IsUbuntu() {
		return false
	}
	switch i.VersionID {
	case UbuntuJammy, UbuntuNoble, UbuntuResolute:
		return true
	}
	return false
}

// IsAlmaLinux8Or9Or10 reports true for AlmaLinux 8, 9, or 10.
func (i Info) IsAlmaLinux8Or9Or10() bool {
	return i.ID == "almalinux" && i.IsSupportedEL()
}

// IsRL9Family reports true for Rocky Linux 8, 9, or 10 (the explicitly
// supported Enterprise Linux releases beyond AlmaLinux).
func (i Info) IsRocky8Or9Or10() bool {
	return i.ID == "rocky" && i.IsSupportedEL()
}

// IsSupportedEL reports true for any Enterprise Linux host on a major version
// this project targets (8, 9, or 10), regardless of the exact rebuild.
func (i Info) IsSupportedEL() bool {
	return i.IsRHEL() && (i.MajorVersion == 8 || i.MajorVersion == 9 || i.MajorVersion == 10)
}

// PackageManager returns "dnf" on RHEL-family hosts and "apt" otherwise.
func (i Info) PackageManager() string {
	if i.IsRHEL() {
		return "dnf"
	}
	return "apt"
}

// SSHServiceName returns the canonical OpenSSH systemd unit basename.
func (i Info) SSHServiceName() string {
	if i.IsRHEL() {
		return "sshd"
	}
	return "ssh"
}

// RedisServiceName returns the cache daemon systemd unit basename:
// "valkey" on RHEL 10+ (AlmaLinux 10), "redis" on RHEL 8/9, "redis-server" on Debian/Ubuntu.
func (i Info) RedisServiceName() string {
	if i.IsRHEL() {
		if i.MajorVersion >= 10 {
			return "valkey"
		}
		return "redis"
	}
	return "redis-server"
}

// FirewallBackend returns "firewalld" on RHEL-family hosts and "ufw" otherwise.
func (i Info) FirewallBackend() string {
	if i.IsRHEL() {
		return "firewalld"
	}
	return "ufw"
}

// MACBackend returns "selinux" on RHEL-family hosts and "apparmor" otherwise.
func (i Info) MACBackend() string {
	if i.IsRHEL() {
		return "selinux"
	}
	return "apparmor"
}

// NologinShell returns the canonical nologin binary path.
func (i Info) NologinShell() string {
	if i.IsRHEL() {
		return "/sbin/nologin"
	}
	return "/usr/sbin/nologin"
}

// PowerDNSRecursorForwardsPath returns the recursor forwarder file path.
func (i Info) PowerDNSRecursorForwardsPath() string {
	if i.IsRHEL() {
		return "/etc/pdns-recursor/recursor.forwards"
	}
	return "/etc/powerdns/recursor.forwards"
}

// RemiShortVersion converts "8.4" into Remi's SCL token "84".
func RemiShortVersion(version string) string {
	return strings.ReplaceAll(strings.TrimSpace(version), ".", "")
}

// ExpandRemiShortVersion converts Remi's SCL token "84" (or "74") into "8.4".
func ExpandRemiShortVersion(short string) (string, bool) {
	s := strings.TrimSpace(short)
	if len(s) < 2 {
		return "", false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return "", false
		}
	}
	return s[:1] + "." + s[1:], true
}

// PHPPoolDir returns the FPM pool directory for the given X.Y PHP version.
func (i Info) PHPPoolDir(version string) string {
	if i.IsRHEL() {
		return filepath.Join("/etc/opt/remi", "php"+RemiShortVersion(version), "php-fpm.d")
	}
	return filepath.Join("/etc/php", version, "fpm", "pool.d")
}

// PHPConfDDir returns the SAPI ini scan directory for the given PHP version and SAPI ("fpm" or "cli").
func (i Info) PHPConfDDir(version, sapi string) string {
	if i.IsRHEL() {
		return filepath.Join("/etc/opt/remi", "php"+RemiShortVersion(version), "php.d")
	}
	return filepath.Join("/etc/php", version, sapi, "conf.d")
}

// PHPIniPath returns the main php.ini path for FPM reads.
func (i Info) PHPIniPath(version string) string {
	if i.IsRHEL() {
		return filepath.Join("/etc/opt/remi", "php"+RemiShortVersion(version), "php.ini")
	}
	return filepath.Join("/etc/php", version, "fpm", "php.ini")
}

// PHPEtcVersionDir returns the top-level /etc tree for a PHP version.
func (i Info) PHPEtcVersionDir(version string) string {
	if i.IsRHEL() {
		return filepath.Join("/etc/opt/remi", "php"+RemiShortVersion(version))
	}
	return filepath.Join("/etc/php", version)
}

// PHPFPMBinary returns the canonical php-fpm binary path for the PHP version.
func (i Info) PHPFPMBinary(version string) string {
	if i.IsRHEL() {
		return filepath.Join("/opt/remi", "php"+RemiShortVersion(version), "root", "usr", "sbin", "php-fpm")
	}
	return "/usr/sbin/php-fpm" + version
}

// PHPCLIBinary returns the canonical CLI php binary path for the PHP version.
func (i Info) PHPCLIBinary(version string) string {
	if i.IsRHEL() {
		return "/usr/bin/php" + RemiShortVersion(version)
	}
	return "/usr/bin/php" + version
}

// PHPFPMServiceUnit returns the systemd service unit name for the distro's global FPM service.
func (i Info) PHPFPMServiceUnit(version string) string {
	if i.IsRHEL() {
		return "php" + RemiShortVersion(version) + "-php-fpm.service"
	}
	return "php" + version + "-fpm.service"
}

// PHPPackagePrefix returns the package name prefix ("php8.4-" on Debian, "php84-php-" on Remi).
func (i Info) PHPPackagePrefix(version string) string {
	if i.IsRHEL() {
		return "php" + RemiShortVersion(version) + "-php-"
	}
	return "php" + version + "-"
}

// PHPCommonPackage returns the -common package name for the PHP version.
func (i Info) PHPCommonPackage(version string) string {
	return i.PHPPackagePrefix(version) + "common"
}

// PHPModuleSearchGlobs returns glob patterns where compiled PHP extension .so modules reside.
func (i Info) PHPModuleSearchGlobs(version, soName string) []string {
	if i.IsRHEL() {
		short := RemiShortVersion(version)
		return []string{
			filepath.Join("/opt/remi", "php"+short, "root", "usr", "lib64", "php", "modules", soName),
			filepath.Join("/usr/lib64/php/modules", soName),
		}
	}
	return []string{filepath.Join("/usr/lib/php/*", soName)}
}

// ExtractPHPVersionFromPoolDir parses a pool directory path in either Debian
// (/etc/php/8.4/fpm/pool.d) or Remi (/etc/opt/remi/php84/php-fpm.d) format.
func ExtractPHPVersionFromPoolDir(poolDir string) (string, bool) {
	cleaned := filepath.Clean(poolDir)
	parts := strings.Split(cleaned, string(filepath.Separator))
	// Debian: ["", "etc", "php", "8.4", "fpm", "pool.d"]
	if len(parts) == 6 && parts[1] == "etc" && parts[2] == "php" && parts[4] == "fpm" && parts[5] == "pool.d" {
		if strings.Contains(parts[3], ".") {
			return parts[3], true
		}
	}
	// Remi EL: ["", "etc", "opt", "remi", "php84", "php-fpm.d"]
	if len(parts) == 6 && parts[1] == "etc" && parts[2] == "opt" && parts[3] == "remi" && strings.HasPrefix(parts[4], "php") && parts[5] == "php-fpm.d" {
		return ExpandRemiShortVersion(strings.TrimPrefix(parts[4], "php"))
	}
	return "", false
}
