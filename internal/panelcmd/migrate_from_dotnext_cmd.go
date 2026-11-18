// SPDX-License-Identifier: MIT

package panelcmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type dotnextUpgradePlan struct {
	RootPrefix string
	DryRun     bool
	RunCmd     func(name string, args ...string) error
	Out        io.Writer
}

// The legacy source is a pre-dotNext panel. Its paths, unit names and env
// prefixes stay spelled as that product shipped them: these are what we must
// detect and rewrite, so they are NOT brand-renamed.
var legacyJabaliServices = []string{
	"jabali-panel.service",
	"jabali-agent.service",
	"jabali-kratos.service",
	"jabali-mailrelay.service",
	"jabali-mailhook.service",
	"jabali-stalwart.service",
	"jabali-webmail.service",
	"jabali-webdav-auth.service",
	"jabali-sso-reaper.timer",
	"jabali-retention-sweep.timer",
	"jabali-cache-doctor.timer",
	"jabali-disk-maintenance.timer",
	"jabali-aide-check.timer",
	"jabali-crowdsec-hub-refresh.timer",
	"jabali-security-feed.timer",
	"jabali-spam-rules-update.timer",
	"jabali-backup-retention.timer",
}

var legacyDirMoves = [][2]string{
	{"/var/lib/jabali-panel", "/var/lib/dotnext"},
	{"/var/lib/jabali", "/var/lib/dotnext"},
	{"/var/lib/jabali-backups", "/var/lib/dotnext-backups"},
	{"/var/lib/jabali-uploads", "/var/lib/dotnext-uploads"},
	{"/var/lib/jabali-agent-staging", "/var/lib/dotnext-agent-staging"},
	{"/var/lib/jabali-ftp-jails", "/var/lib/dotnext-ftp-jails"},
	{"/var/lib/jabali-webmail", "/var/lib/dotnext-webmail"},
	{"/var/log/jabali", "/var/log/dotnext"},
}

var legacyConfigSources = []string{
	"/etc/jabali",
	"/etc/jabali-panel",
}

var legacyEnvReplacements = [][2]string{
	{"/etc/jabali-panel", "/etc/dotnext"},
	{"/etc/jabali", "/etc/dotnext"},
	{"/run/jabali-panel", "/run/dotnext-panel"},
	{"/run/jabali-kratos", "/run/dotnext-kratos"},
	{"/run/jabali-webdav", "/run/dotnext-webdav"},
	{"/run/jabali-mailrelay", "/run/dotnext-mailrelay"},
	{"/run/jabali", "/run/dotnext"},
	{"/var/lib/jabali-panel", "/var/lib/dotnext"},
	{"/var/lib/jabali-backups", "/var/lib/dotnext-backups"},
	{"/var/lib/jabali-uploads", "/var/lib/dotnext-uploads"},
	{"/var/lib/jabali-agent-staging", "/var/lib/dotnext-agent-staging"},
	{"/var/lib/jabali-ftp-jails", "/var/lib/dotnext-ftp-jails"},
	{"/var/lib/jabali-webmail", "/var/lib/dotnext-webmail"},
	{"/var/lib/jabali", "/var/lib/dotnext"},
	{"/opt/jabali-panel", "/opt/dotnext"},
	{"/opt/jabali-webmail", "/opt/dotnext-webmail"},
	{"JABALI_", "DOTNEXT_"},
}

func newMigrateFromDotNextCmd() *cobra.Command {
	var dryRun bool
	var apply bool

	cmd := &cobra.Command{
		Use:   "migrate-from-dotnext",
		Short: "Migrate an existing legacy DotNext Panel host layout to dotNext",
		Long: "Detects legacy /etc/dotnext*, /var/lib/dotnext*, dotnext-* systemd units, " +
			"and /usr/local/bin/dotnext* binaries, and cleanly migrates configuration " +
			"and state to pure dotNext paths without leaving runtime aliases.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !apply && !dryRun {
				dryRun = true
			}
			if apply && dryRun {
				return fmt.Errorf("--apply and --dry-run are mutually exclusive")
			}
			plan := dotnextUpgradePlan{
				RootPrefix: "/",
				DryRun:     dryRun,
				RunCmd: func(name string, args ...string) error {
					c := exec.Command(name, args...)
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					return c.Run()
				},
				Out: cmd.OutOrStdout(),
			}
			return runDotNextUpgradeMigration(plan)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview planned migrations without mutating files (default unless --apply is passed)")
	cmd.Flags().BoolVar(&apply, "apply", false, "perform the in-place migration to dotNext paths and units")
	return cmd
}

func runDotNextUpgradeMigration(plan dotnextUpgradePlan) error {
	out := plan.Out
	if out == nil {
		out = io.Discard
	}
	root := plan.RootPrefix
	if root == "" {
		root = "/"
	}
	joinRoot := func(abs string) string {
		return filepath.Join(root, strings.TrimPrefix(abs, "/"))
	}

	foundLegacy := false
	for _, src := range legacyConfigSources {
		if st, err := os.Stat(joinRoot(src)); err == nil && st.IsDir() {
			foundLegacy = true
		}
	}
	for _, pair := range legacyDirMoves {
		if st, err := os.Stat(joinRoot(pair[0])); err == nil && st.IsDir() {
			foundLegacy = true
		}
	}
	for _, unit := range legacyJabaliServices {
		if _, err := os.Stat(joinRoot(filepath.Join("/etc/systemd/system", unit))); err == nil {
			foundLegacy = true
		}
	}

	if !foundLegacy {
		fmt.Fprintln(out, "No legacy DotNext installation detected; nothing to migrate.")
		return nil
	}

	for _, unit := range legacyJabaliServices {
		unitPath := joinRoot(filepath.Join("/etc/systemd/system", unit))
		if _, err := os.Stat(unitPath); err != nil {
			continue
		}
		fmt.Fprintf(out, "stop/disable/remove legacy unit %s\n", unit)
		if !plan.DryRun {
			if plan.RunCmd != nil {
				_ = plan.RunCmd("systemctl", "stop", unit)
				_ = plan.RunCmd("systemctl", "disable", unit)
			}
			_ = os.Remove(unitPath)
		}
	}

	dstCfg := joinRoot("/etc/dotnext")
	if !plan.DryRun {
		if err := os.MkdirAll(dstCfg, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dstCfg, err)
		}
	}
	for _, src := range legacyConfigSources {
		srcPath := joinRoot(src)
		if st, err := os.Stat(srcPath); err != nil || !st.IsDir() {
			continue
		}
		fmt.Fprintf(out, "migrate config directory %s -> /etc/dotnext\n", src)
		if !plan.DryRun {
			if err := copyMissingTree(srcPath, dstCfg); err != nil {
				return fmt.Errorf("copy %s -> %s: %w", srcPath, dstCfg, err)
			}
		}
	}

	for _, name := range []string{"panel.env", "config.toml", "pdns.env", "kratos.yml", "bulwark.env"} {
		cfgFile := filepath.Join(dstCfg, name)
		raw, err := os.ReadFile(cfgFile)
		if err != nil {
			continue
		}
		updated := rewriteLegacyConfigText(string(raw))
		if updated != string(raw) {
			fmt.Fprintf(out, "rewrite legacy references inside /etc/dotnext/%s\n", name)
			if !plan.DryRun {
				info, _ := os.Stat(cfgFile)
				mode := fs.FileMode(0o640)
				if info != nil {
					mode = info.Mode()
				}
				if err := os.WriteFile(cfgFile, []byte(updated), mode); err != nil {
					return fmt.Errorf("write %s: %w", cfgFile, err)
				}
			}
		}
	}

	for _, pair := range legacyDirMoves {
		oldDir := joinRoot(pair[0])
		newDir := joinRoot(pair[1])
		if st, err := os.Stat(oldDir); err != nil || !st.IsDir() {
			continue
		}
		if _, err := os.Stat(newDir); err == nil {
			fmt.Fprintf(out, "skip move %s -> %s (destination already exists)\n", pair[0], pair[1])
			continue
		}
		fmt.Fprintf(out, "move state directory %s -> %s\n", pair[0], pair[1])
		if !plan.DryRun {
			if err := os.MkdirAll(filepath.Dir(newDir), 0o755); err != nil {
				return fmt.Errorf("mkdir parent %s: %w", newDir, err)
			}
			if err := os.Rename(oldDir, newDir); err != nil {
				return fmt.Errorf("rename %s -> %s: %w", oldDir, newDir, err)
			}
		}
	}

	// Legacy binaries carry the same brand name the legacy units and paths do;
	// they are removed so no stale alias survives the upgrade.
	for _, legacyBin := range []string{
		"/usr/local/bin/jabali",
		"/usr/local/bin/jabali-panel",
		"/usr/local/bin/jabali-agent",
		"/usr/local/bin/jabali-webdav",
		"/usr/local/bin/jabali-ssh-shell",
		"/usr/local/bin/jabali-mailhook",
		"/usr/local/libexec/jabali/jabali-sendmail",
	} {
		target := joinRoot(legacyBin)
		if _, err := os.Lstat(target); err == nil {
			fmt.Fprintf(out, "remove legacy binary %s\n", legacyBin)
			if !plan.DryRun {
				_ = os.Remove(target)
			}
		}
	}

	if !plan.DryRun && plan.RunCmd != nil {
		_ = plan.RunCmd("systemctl", "daemon-reload")
	}
	return nil
}

func rewriteLegacyConfigText(in string) string {
	out := in
	for _, pair := range legacyEnvReplacements {
		out = strings.ReplaceAll(out, pair[0], pair[1])
	}
	return out
}

func copyMissingTree(srcRoot, dstRoot string) error {
	return filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(dstRoot, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(dst, info.Mode().Perm())
		}
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode().Perm())
	})
}
