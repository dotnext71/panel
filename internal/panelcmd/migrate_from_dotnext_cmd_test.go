// SPDX-License-Identifier: MIT

package panelcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDotNextUpgradeMigration_NoLegacyIsNoop(t *testing.T) {
	var out bytes.Buffer
	err := runDotNextUpgradeMigration(dotnextUpgradePlan{
		RootPrefix: t.TempDir(),
		DryRun:     false,
		Out:        &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "No legacy DotNext installation detected") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunDotNextUpgradeMigration_DryRunDoesNotMutate(t *testing.T) {
	// The legacy config directory is what marks the host as a pre-dotNext
	// install; a plain /etc/dotnext is the CURRENT layout, not a legacy one.
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "etc", legacyPrefix))
	mustWriteFile(t, filepath.Join(root, "etc", legacyPrefix, "panel.env"), "PANEL_ADDR=unix:/run/"+legacyPrefix+"-panel/api.sock\n"+strings.ToUpper(legacyPrefix)+"_ENV=prod\n")

	var cmds [][]string
	var out bytes.Buffer
	err := runDotNextUpgradeMigration(dotnextUpgradePlan{
		RootPrefix: root,
		DryRun:     true,
		RunCmd: func(name string, args ...string) error {
			cmds = append(cmds, append([]string{name}, args...))
			return nil
		},
		Out: &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 0 {
		t.Fatalf("dry-run invoked commands: %v", cmds)
	}
	if _, err := os.Stat(filepath.Join(root, "etc/dotnext/panel.env")); err == nil {
		t.Fatal("dry-run created /etc/dotnext/panel.env")
	}
}

func TestRunDotNextUpgradeMigration_ApplyMigratesLayoutAndConfig(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "etc", legacyPrefix))
	mustMkdirAll(t, filepath.Join(root, "etc/systemd/system"))
	mustMkdirAll(t, filepath.Join(root, "var/lib", legacyPrefix+"-backups", "repo"))
	mustMkdirAll(t, filepath.Join(root, "usr/local/bin"))

	mustWriteFile(t, filepath.Join(root, "etc", legacyPrefix, "panel.env"),
		"PANEL_ADDR=unix:/run/"+legacyPrefix+"-panel/api.sock\n"+
			"AGENT_SOCKET=/run/"+legacyPrefix+"/agent.sock\n"+
			strings.ToUpper(legacyPrefix)+"_ADMIN=yes\n")
	mustWriteFile(t, filepath.Join(root, "etc/systemd/system", legacyPrefix+"-panel.service"), "[Unit]\n")
	mustWriteFile(t, filepath.Join(root, "usr/local/bin", legacyPrefix), "legacy")

	var cmds []string
	var out bytes.Buffer
	err := runDotNextUpgradeMigration(dotnextUpgradePlan{
		RootPrefix: root,
		DryRun:     false,
		RunCmd: func(name string, args ...string) error {
			cmds = append(cmds, name+" "+strings.Join(args, " "))
			return nil
		},
		Out: &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotEnv, err := os.ReadFile(filepath.Join(root, "etc/dotnext/panel.env"))
	if err != nil {
		t.Fatalf("read migrated panel.env: %v", err)
	}
	wantEnv := "PANEL_ADDR=unix:/run/dotnext-panel/api.sock\nAGENT_SOCKET=/run/dotnext/agent.sock\nDOTNEXT_ADMIN=yes\n"
	if string(gotEnv) != wantEnv {
		t.Fatalf("migrated env = %q, want %q", string(gotEnv), wantEnv)
	}

	if _, err := os.Stat(filepath.Join(root, "var/lib/dotnext-backups/repo")); err != nil {
		t.Fatalf("expected /var/lib/dotnext-backups/repo: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "etc/systemd/system", legacyPrefix+"-panel.service")); !os.IsNotExist(err) {
		t.Fatalf("expected legacy unit removal, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "usr/local/bin", legacyPrefix)); !os.IsNotExist(err) {
		t.Fatalf("expected legacy binary removal, stat err = %v", err)
	}
	if len(cmds) < 3 {
		t.Fatalf("expected systemctl stop/disable/daemon-reload calls, got %v", cmds)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
