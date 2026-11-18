package panelcmd

import (
	"runtime"
	"testing"
)

// GH #727: on a Custom / no-mail install the dotnext-mail group never exists, so
// the mail convergence buildSteps (stalwart binary, spam-filter rules,
// mailhook, bulwark) must be skipped on update — otherwise
// _install_spam_rules dies at `install -g dotnext-mail`.
func TestGroupExists_GH727(t *testing.T) {
	rootGroup := "root"
	if runtime.GOOS == "darwin" {
		rootGroup = "wheel"
	}
	if !groupExists(rootGroup) {
		t.Errorf("expected the %s group to exist", rootGroup)
	}
	// A group name that cannot exist -> false: this is the no-mail host case
	// that used to let the mail buildSteps run and crash.
	if groupExists("dotnext-no-such-group-gh727-xyz") {
		t.Error("expected a bogus group to be reported absent")
	}
}

func TestMailModuleInstalled_TracksDotNextMailGroup_GH727(t *testing.T) {
	if mailModuleInstalled() != groupExists("dotnext-mail") {
		t.Error("mailModuleInstalled must track presence of the dotnext-mail group")
	}
}
