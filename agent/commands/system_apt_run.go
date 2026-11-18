package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dotnext-panel/dotnext/agentwire"
	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// systemAptRunResponse mirrors systemUpdateRunResponse — the apt/dnf path
// detaches into its own transient cgroup so a libc/openssh upgrade that
// triggers a daemon restart doesn't kill the upgrade process running it.
type systemAptRunResponse struct {
	Unit      string `json:"unit"`
	StartedAt string `json:"started_at"`
}

const aptUnitName = "dotnext-apt-oneshot.service"

func systemAptRunHandler(ctx context.Context, _ json.RawMessage) (any, error) {
	_ = execCommandContext(ctx, "systemctl", "reset-failed", aptUnitName).Run()

	startedAt := time.Now().UTC()
	args := []string{
		"--unit=" + aptUnitName,
		"--no-block",
		// No --collect: failed unit lingers (reconciler reads it as failed),
		// success is GC'd to inactive. reset-failed at the next run clears it.
		"--setenv=DEBIAN_FRONTEND=noninteractive",
		"--setenv=LC_ALL=C",
		"apt-get",
		"-y",
		"-o=Dpkg::Options::=--force-confdef",
		"-o=Dpkg::Options::=--force-confold",
		"-o=DPkg::Lock::Timeout=120",
		"dist-upgrade",
	}
	if osinfo.Detect().IsRHEL() {
		args = []string{
			"--unit=" + aptUnitName,
			"--no-block",
			"--setenv=LC_ALL=C",
			"dnf", "-y", "upgrade",
		}
	}
	cmd := execCommandContext(ctx, "systemd-run", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInternal,
			Message: fmt.Sprintf("systemd-run: %v: %s", err, string(out)),
		}
	}
	return systemAptRunResponse{
		Unit:      aptUnitName,
		StartedAt: startedAt.Format(time.RFC3339Nano),
	}, nil
}

func init() {
	Default.Register("system.apt_run", systemAptRunHandler)
}
