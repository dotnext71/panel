package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractVersionFromPath_DebianAndRemi(t *testing.T) {
	if v, err := extractVersionFromPath("/etc/php/8.3/fpm/pool.d"); err != nil || v != "8.3" {
		t.Fatalf("debian path = (%q, %v), want 8.3", v, err)
	}
	if v, err := extractVersionFromPath("/etc/opt/remi/php84/php-fpm.d"); err != nil || v != "8.4" {
		t.Fatalf("remi path = (%q, %v), want 8.4", v, err)
	}
}

func TestToggleRemiPHPIni(t *testing.T) {
	dir := t.TempDir()
	ini := filepath.Join(dir, "40-igbinary.ini")
	if err := os.WriteFile(ini, []byte("extension=igbinary.so\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := toggleRemiPHPIni(dir, "igbinary", false); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if _, err := os.Stat(ini + ".disabled"); err != nil {
		t.Fatalf("expected .disabled file: %v", err)
	}
	if _, err := toggleRemiPHPIni(dir, "igbinary", true); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if _, err := os.Stat(ini); err != nil {
		t.Fatalf("expected re-enabled file: %v", err)
	}
}
