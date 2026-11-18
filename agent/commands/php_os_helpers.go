package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// phpPoolDirForVersion returns the FPM pool directory for the installed PHP
// version on the current OS family (Debian /etc/php/<v>/fpm/pool.d vs Remi EL
// /etc/opt/remi/php<nodot>/php-fpm.d). If an override exists on disk in one of
// the two trees, it wins automatically.
func phpPoolDirForVersion(version string) string {
	info := osinfo.Detect()
	primary := info.PHPPoolDir(version)
	if st, err := os.Stat(primary); err == nil && st.IsDir() {
		return primary
	}
	deb := filepath.Join("/etc/php", version, "fpm", "pool.d")
	remi := filepath.Join("/etc/opt/remi", "php"+osinfo.RemiShortVersion(version), "php-fpm.d")
	if info.IsRHEL() {
		if st, err := os.Stat(deb); err == nil && st.IsDir() {
			return deb
		}
		return remi
	}
	if st, err := os.Stat(remi); err == nil && st.IsDir() {
		return remi
	}
	return deb
}

// phpFPMBinaryForVersion returns the php-fpm binary path for the given version,
// probing the host's preferred layout first and falling back to the alternate tree.
func phpFPMBinaryForVersion(version string) string {
	info := osinfo.Detect()
	deb := "/usr/sbin/php-fpm" + version
	remi := filepath.Join("/opt/remi", "php"+osinfo.RemiShortVersion(version), "root", "usr", "sbin", "php-fpm")
	if info.IsRHEL() {
		if _, err := os.Stat(remi); err == nil {
			return remi
		}
		if _, err := os.Stat(deb); err == nil {
			return deb
		}
		return remi
	}
	if _, err := os.Stat(deb); err == nil {
		return deb
	}
	if _, err := os.Stat(remi); err == nil {
		return remi
	}
	return deb
}

// phpCLIBinaryForVersion returns the CLI php binary path for the given version.
func phpCLIBinaryForVersion(version string) string {
	info := osinfo.Detect()
	deb := "/usr/bin/php" + version
	remi := "/usr/bin/php" + osinfo.RemiShortVersion(version)
	remiRoot := filepath.Join("/opt/remi", "php"+osinfo.RemiShortVersion(version), "root", "usr", "bin", "php")
	if info.IsRHEL() {
		for _, c := range []string{remi, remiRoot, deb} {
			if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
				return c
			}
		}
		return remi
	}
	for _, c := range []string{deb, remi} {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return deb
}

// globPoolFilesForSlug finds per-slug FPM pool files across both Debian and
// Remi EL layouts while respecting custom phpEtcRoot overrides in tests.
func globPoolFilesForSlug(slug string) ([]string, error) {
	debPattern := filepath.Join(phpEtcRoot, "*", "fpm", "pool.d", "dotnext-"+slug+".conf")
	matches, err := filepath.Glob(debPattern)
	if err != nil {
		return nil, err
	}
	if phpEtcRoot == "/etc/php" {
		remiPattern := filepath.Join("/etc/opt/remi", "php*", "php-fpm.d", "dotnext-"+slug+".conf")
		if remiMatches, rErr := filepath.Glob(remiPattern); rErr == nil {
			matches = append(matches, remiMatches...)
		}
	}
	return matches, nil
}

// globPoolDirs returns every installed FPM pool directory on disk across
// Debian (/etc/php/<v>/fpm/pool.d) and Remi (/etc/opt/remi/php<nodot>/php-fpm.d).
func globPoolDirs() ([]string, error) {
	matches, err := filepath.Glob("/etc/php/*/fpm/pool.d")
	if err != nil {
		return nil, err
	}
	if remiMatches, rErr := filepath.Glob("/etc/opt/remi/php*/php-fpm.d"); rErr == nil {
		matches = append(matches, remiMatches...)
	}
	return matches, nil
}

// findPHPModuleSo checks whether `<ext>.so` exists in Debian or Remi module directories.
func findPHPModuleSo(version, so string) bool {
	info := osinfo.Detect()
	for _, pattern := range info.PHPModuleSearchGlobs(version, so) {
		if m, _ := filepath.Glob(pattern); len(m) > 0 {
			return true
		}
	}
	return false
}

// phpCLIBinaryName returns the best executable binary name/path for a user's CLI PHP pin.
func phpCLIBinaryName(version string) string {
	if _, err := exec.LookPath("php" + version); err == nil {
		return "php" + version
	}
	short := "php" + osinfo.RemiShortVersion(version)
	if _, err := exec.LookPath(short); err == nil {
		return short
	}
	return phpCLIBinaryForVersion(version)
}
