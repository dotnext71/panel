package osinfo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOSRelease_AlmaLinuxVersions(t *testing.T) {
	cases := []struct {
		name        string
		raw         string
		wantID      string
		wantMajor   int
		wantRHEL    bool
		wantAlma    bool
		wantRedis   string
		wantSSH     string
		wantFW      string
		wantMAC     string
		wantPkgMgr  string
		wantFPMBin  string
		wantPoolDir string
		wantIni     string
	}{
		{
			name: "AlmaLinux 8",
			raw: `ID="almalinux"
ID_LIKE="rhel centos fedora"
VERSION_ID="8.10"
PRETTY_NAME="AlmaLinux 8.10 (Cerulean Leopard)"`,
			wantID:      "almalinux",
			wantMajor:   8,
			wantRHEL:    true,
			wantAlma:    true,
			wantRedis:   "redis",
			wantSSH:     "sshd",
			wantFW:      "firewalld",
			wantMAC:     "selinux",
			wantPkgMgr:  "dnf",
			wantFPMBin:  "/opt/remi/php84/root/usr/sbin/php-fpm",
			wantPoolDir: "/etc/opt/remi/php84/php-fpm.d",
			wantIni:     "/etc/opt/remi/php84/php.ini",
		},
		{
			name: "AlmaLinux 9",
			raw: `ID="almalinux"
ID_LIKE="rhel centos fedora"
VERSION_ID="9.5"
PRETTY_NAME="AlmaLinux 9.5 (Teal Serval)"`,
			wantID:      "almalinux",
			wantMajor:   9,
			wantRHEL:    true,
			wantAlma:    true,
			wantRedis:   "redis",
			wantSSH:     "sshd",
			wantFW:      "firewalld",
			wantMAC:     "selinux",
			wantPkgMgr:  "dnf",
			wantFPMBin:  "/opt/remi/php84/root/usr/sbin/php-fpm",
			wantPoolDir: "/etc/opt/remi/php84/php-fpm.d",
			wantIni:     "/etc/opt/remi/php84/php.ini",
		},
		{
			name: "AlmaLinux 10",
			raw: `ID="almalinux"
ID_LIKE="rhel centos fedora"
VERSION_ID="10.0"
PRETTY_NAME="AlmaLinux 10.0 (Purple Lion)"`,
			wantID:      "almalinux",
			wantMajor:   10,
			wantRHEL:    true,
			wantAlma:    true,
			wantRedis:   "valkey",
			wantSSH:     "sshd",
			wantFW:      "firewalld",
			wantMAC:     "selinux",
			wantPkgMgr:  "dnf",
			wantFPMBin:  "/opt/remi/php84/root/usr/sbin/php-fpm",
			wantPoolDir: "/etc/opt/remi/php84/php-fpm.d",
			wantIni:     "/etc/opt/remi/php84/php.ini",
		},
		{
			name: "Debian 13",
			raw: `ID=debian
VERSION_ID="13"
VERSION_CODENAME=trixie
PRETTY_NAME="Debian GNU/Linux 13 (trixie)"`,
			wantID:      "debian",
			wantMajor:   13,
			wantRHEL:    false,
			wantAlma:    false,
			wantRedis:   "redis-server",
			wantSSH:     "ssh",
			wantFW:      "ufw",
			wantMAC:     "apparmor",
			wantPkgMgr:  "apt",
			wantFPMBin:  "/usr/sbin/php-fpm8.4",
			wantPoolDir: "/etc/php/8.4/fpm/pool.d",
			wantIni:     "/etc/php/8.4/fpm/php.ini",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseOSRelease(tc.raw)
			if got.ID != tc.wantID {
				t.Fatalf("ID = %q, want %q", got.ID, tc.wantID)
			}
			if got.MajorVersion != tc.wantMajor {
				t.Fatalf("MajorVersion = %d, want %d", got.MajorVersion, tc.wantMajor)
			}
			if got.IsRHEL() != tc.wantRHEL {
				t.Fatalf("IsRHEL() = %v, want %v", got.IsRHEL(), tc.wantRHEL)
			}
			if got.IsAlmaLinux8Or9Or10() != tc.wantAlma {
				t.Fatalf("IsAlmaLinux8Or9Or10() = %v, want %v", got.IsAlmaLinux8Or9Or10(), tc.wantAlma)
			}
			if got.RedisServiceName() != tc.wantRedis {
				t.Fatalf("RedisServiceName() = %q, want %q", got.RedisServiceName(), tc.wantRedis)
			}
			if got.SSHServiceName() != tc.wantSSH {
				t.Fatalf("SSHServiceName() = %q, want %q", got.SSHServiceName(), tc.wantSSH)
			}
			if got.FirewallBackend() != tc.wantFW {
				t.Fatalf("FirewallBackend() = %q, want %q", got.FirewallBackend(), tc.wantFW)
			}
			if got.MACBackend() != tc.wantMAC {
				t.Fatalf("MACBackend() = %q, want %q", got.MACBackend(), tc.wantMAC)
			}
			if got.PackageManager() != tc.wantPkgMgr {
				t.Fatalf("PackageManager() = %q, want %q", got.PackageManager(), tc.wantPkgMgr)
			}
			if got.PHPFPMBinary("8.4") != tc.wantFPMBin {
				t.Fatalf("PHPFPMBinary(8.4) = %q, want %q", got.PHPFPMBinary("8.4"), tc.wantFPMBin)
			}
			if got.PHPPoolDir("8.4") != tc.wantPoolDir {
				t.Fatalf("PHPPoolDir(8.4) = %q, want %q", got.PHPPoolDir("8.4"), tc.wantPoolDir)
			}
			if got.PHPIniPath("8.4") != tc.wantIni {
				t.Fatalf("PHPIniPath(8.4) = %q, want %q", got.PHPIniPath("8.4"), tc.wantIni)
			}
		})
	}
}

// Rocky Linux is a binary-compatible rebuild: everything that keys off the
// RHEL family must behave identically to AlmaLinux, and the family-wide
// predicate must accept both.
// Ubuntu releases supported by the panel: the three LTS releases. 26.04
// (Resolute Raccoon) is the version this change adds; interim releases such as
// 25.10 are intentionally unsupported.
func TestParseOSRelease_UbuntuReleases(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		versionID string
		supported bool
	}{
		{"22.04 jammy", "ID=ubuntu\nVERSION_ID=\"22.04\"\nVERSION_CODENAME=jammy", "22.04", true},
		{"24.04 noble", "ID=ubuntu\nVERSION_ID=\"24.04\"\nVERSION_CODENAME=noble", "24.04", true},
		// 25.10 is EOL (July 2026) and has no PHP suite in any repo: unsupported.
		{"25.10 questing", "ID=ubuntu\nVERSION_ID=\"25.10\"\nVERSION_CODENAME=questing", "25.10", false},
		{"26.04 resolute", "ID=ubuntu\nVERSION_ID=\"26.04\"\nVERSION_CODENAME=resolute", "26.04", true},
		{"24.10 oracular", "ID=ubuntu\nVERSION_ID=\"24.10\"", "24.10", false},
		{"18.04 bionic", "ID=ubuntu\nVERSION_ID=\"18.04\"", "18.04", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			info := ParseOSRelease(tc.raw)
			if !info.IsUbuntu() {
				t.Fatal("IsUbuntu() = false")
			}
			if !info.IsDebian() {
				t.Fatal("Ubuntu must belong to the Debian family")
			}
			if info.IsSupportedUbuntu() != tc.supported {
				t.Fatalf("IsSupportedUbuntu() = %v, want %v", info.IsSupportedUbuntu(), tc.supported)
			}
			if !info.IsUbuntuRelease(tc.versionID) {
				t.Fatalf("IsUbuntuRelease(%q) = false", tc.versionID)
			}
		})
	}
}

// Debian must never be mistaken for Ubuntu: the PHP/nginx/Docker sources differ
// per release, and an "ubuntu" answer on Debian would select a nonexistent suite.
func TestUbuntuPredicates_DoNotMatchDebian(t *testing.T) {
	for _, raw := range []string{
		"ID=debian\nVERSION_ID=\"13\"\nVERSION_CODENAME=trixie",
		"ID=debian\nVERSION_ID=\"12\"\nVERSION_CODENAME=bookworm",
	} {
		info := ParseOSRelease(raw)
		if info.IsUbuntu() || info.IsSupportedUbuntu() || info.IsUbuntuRelease("24.04") {
			t.Fatalf("Debian (%s) matched a Ubuntu predicate: %+v", info.VersionID, info)
		}
		if !info.IsDebian() {
			t.Fatalf("Debian %s not IsDebian()", info.VersionID)
		}
	}
	// …and the reverse: Ubuntu must not satisfy the EL predicates.
	u := ParseOSRelease("ID=ubuntu\nVERSION_ID=\"24.04\"")
	if u.IsRHEL() || u.IsSupportedEL() || u.IsAlmaLinux8Or9Or10() || u.IsRocky8Or9Or10() {
		t.Fatalf("Ubuntu matched an EL predicate: %+v", u)
	}
}

func TestParseOSRelease_RockyParityWithAlma(t *testing.T) {
	for _, major := range []string{"8.10", "9.5", "10.0"} {
		t.Run("Rocky_"+major, func(t *testing.T) {
			rocky := ParseOSRelease("ID=\"rocky\"\nID_LIKE=\"rhel centos fedora\"\nVERSION_ID=\"" + major + "\"")
			alma := ParseOSRelease("ID=\"almalinux\"\nID_LIKE=\"rhel centos fedora\"\nVERSION_ID=\"" + major + "\"")

			if !rocky.IsRHEL() {
				t.Fatal("Rocky must be RHEL family")
			}
			if !rocky.IsSupportedEL() || !alma.IsSupportedEL() {
				t.Fatalf("IsSupportedEL: rocky=%v alma=%v, want both true", rocky.IsSupportedEL(), alma.IsSupportedEL())
			}
			if !rocky.IsRocky8Or9Or10() || rocky.IsAlmaLinux8Or9Or10() {
				t.Fatal("Rocky must match IsRocky8Or9Or10 and not IsAlmaLinux8Or9Or10")
			}
			if alma.IsRocky8Or9Or10() || !alma.IsAlmaLinux8Or9Or10() {
				t.Fatal("Alma must match IsAlmaLinux8Or9Or10 and not IsRocky8Or9Or10")
			}
			// Everything both rebuilds must agree on.
			pairs := map[string][2]string{
				"package manager": {rocky.PackageManager(), alma.PackageManager()},
				"ssh unit":        {rocky.SSHServiceName(), alma.SSHServiceName()},
				"redis unit":      {rocky.RedisServiceName(), alma.RedisServiceName()},
				"firewall":        {rocky.FirewallBackend(), alma.FirewallBackend()},
				"MAC":             {rocky.MACBackend(), alma.MACBackend()},
				"nologin":         {rocky.NologinShell(), alma.NologinShell()},
				"pdns forwards":   {rocky.PowerDNSRecursorForwardsPath(), alma.PowerDNSRecursorForwardsPath()},
				"php pool dir":    {rocky.PHPPoolDir("8.4"), alma.PHPPoolDir("8.4")},
				"php fpm binary":  {rocky.PHPFPMBinary("8.4"), alma.PHPFPMBinary("8.4")},
			}
			for name, v := range pairs {
				if v[0] != v[1] {
					t.Errorf("%s: rocky=%q alma=%q", name, v[0], v[1])
				}
			}
		})
	}
}

func TestExtractPHPVersionFromPoolDir(t *testing.T) {
	if v, ok := ExtractPHPVersionFromPoolDir("/etc/php/8.3/fpm/pool.d"); !ok || v != "8.3" {
		t.Fatalf("debian pool dir = (%q, %v), want (8.3, true)", v, ok)
	}
	if v, ok := ExtractPHPVersionFromPoolDir("/etc/opt/remi/php84/php-fpm.d"); !ok || v != "8.4" {
		t.Fatalf("remi pool dir = (%q, %v), want (8.4, true)", v, ok)
	}
	if v, ok := ExtractPHPVersionFromPoolDir("/etc/opt/remi/php74/php-fpm.d/"); !ok || v != "7.4" {
		t.Fatalf("remi pool dir 7.4 = (%q, %v), want (7.4, true)", v, ok)
	}
	if _, ok := ExtractPHPVersionFromPoolDir("/etc/nginx/sites-enabled"); ok {
		t.Fatal("expected false for unrelated dir")
	}
}

func TestDetectFrom_File(t *testing.T) {
	p := filepath.Join(t.TempDir(), "os-release")
	if err := os.WriteFile(p, []byte("ID=almalinux\nVERSION_ID=\"9.4\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	info := DetectFrom(p)
	if !info.IsAlmaLinux8Or9Or10() || info.MajorVersion != 9 {
		t.Fatalf("unexpected info: %+v", info)
	}
}
