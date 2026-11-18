package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dotnext-panel/dotnext/agentwire"
	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// phpVersionUninstallParams is the input shape for php.version.uninstall.
type phpVersionUninstallParams struct {
	Version string `json:"version"`
}

// phpVersionUninstallResponse is the output shape for php.version.uninstall.
type phpVersionUninstallResponse struct {
	Version    string `json:"version"`
	Installed  bool   `json:"installed"`
	FPMRunning bool   `json:"fpm_running"`
}

// purgePackages runs apt-get purge with autoremove for the given packages.
// Mirrors installPackages' env + output capture. Captures stderr on
// non-zero exit so the agent surfaces a useful diagnostic to the UI
// instead of "exit status 100".
func purgePackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}
	var cmd *exec.Cmd
	if osinfo.Detect().IsRHEL() {
		cmd = execCommand("dnf", append([]string{"remove", "-y"}, pkgs...)...)
	} else {
		cmd = execCommand("apt-get", append(
			[]string{"purge", "-y", "--auto-remove"},
			pkgs...,
		)...)
		cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	}
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("package remove failed: %w; output: %s", err, out.String())
	}
	return nil
}

// phpVersionUninstallHandler is the agent handler for php.version.uninstall.
//
// Removes the distro's php<v>-* packages via apt-get purge --auto-remove.
// The pre-install mask symlink (/etc/systemd/system/php<v>-fpm.service →
// /dev/null) stays put — apt won't recreate the unit on a future reinstall
// without it, and a follow-up install handler unmasks before installing.
//
// Caller (panel-api) is responsible for refusing the call when:
//   - the target is the current default PHP version, or
//   - php-fpm pools still reference it.
//
// Both checks need the panel-api DB; agent does best-effort + reports.
func phpVersionUninstallHandler(_ context.Context, params json.RawMessage) (any, error) {
	if len(params) == 0 {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInvalidArgument,
			Message: "version parameter required",
		}
	}
	var p phpVersionUninstallParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInvalidArgument,
			Message: fmt.Sprintf("failed to parse params: %v", err),
		}
	}
	if !versionRegex.MatchString(p.Version) {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInvalidArgument,
			Message: fmt.Sprintf("invalid version format: %q (expected X.Y)", p.Version),
		}
	}
	if !isVersionSupported(p.Version) {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInvalidArgument,
			Message: fmt.Sprintf("unsupported version: %q", p.Version),
		}
	}

	// Already gone — return idempotent success rather than error.
	if !isInstalledPHPVersion(p.Version) {
		return phpVersionUninstallResponse{
			Version:    p.Version,
			Installed:  false,
			FPMRunning: false,
		}, nil
	}

	info := osinfo.Detect()
	_ = execCommand("systemctl", "stop", info.PHPFPMServiceUnit(p.Version)).Run()

	pkgPattern := info.PHPPackagePrefix(p.Version) + "*"
	if err := purgePackages([]string{pkgPattern}); err != nil {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInternal,
			Message: fmt.Sprintf("remove %s: %v", pkgPattern, err),
		}
	}

	_ = os.RemoveAll(info.PHPEtcVersionDir(p.Version))
	_ = os.RemoveAll(filepath.Join("/etc/php", p.Version))

	// Confirm the runtime binaries are actually gone. If they remain, the
	// purge silently failed (held packages, apt error) — surface a real error
	// instead of reporting a phantom success the UI would trust (GH #187).
	if isInstalledPHPVersion(p.Version) {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInternal,
			Message: fmt.Sprintf("php %s is still installed after purge (binaries present) — packages may be held or apt failed", p.Version),
		}
	}

	return phpVersionUninstallResponse{
		Version:    p.Version,
		Installed:  false,
		FPMRunning: checkFPMRunning(p.Version),
	}, nil
}

func init() {
	Default.Register("php.version.uninstall", phpVersionUninstallHandler)
}
