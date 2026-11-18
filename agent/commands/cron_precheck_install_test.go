package commands

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dotnext-panel/dotnext/internal/cronvalidate"
)

// libexecRefRe finds every /usr/local/libexec/dotnext/<name> path a
// generated unit references in an ExecStart* line.
var libexecRefRe = regexp.MustCompile(`/usr/local/libexec/dotnext/([A-Za-z0-9._-]+)`)

// TestGeneratedUnitLibexecHelpersAreShipped guards the cron-precheck
// 203/EXEC regression: a unit builder must never reference a libexec
// helper that install.sh doesn't ship. For every /usr/local/libexec/
// dotnext/<name> in a rendered cron service unit, assert infrastructure/systemd/
// <name> exists in the repo AND install.sh installs it.
func TestGeneratedUnitLibexecHelpersAreShipped(t *testing.T) {
	cmd, err := cronvalidate.ValidateCommand("php /home/u/public_html/wp-cron.php", []string{"/home/u/public_html"}, "")
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	svc := buildCronServiceContent("job1", "test", []*cronvalidate.Command{cmd}, "u", []string{"/home/u/public_html"})

	refs := libexecRefRe.FindAllStringSubmatch(svc, -1)
	if len(refs) == 0 {
		t.Skip("no libexec helper referenced in the generated unit")
	}

	// repo root = panel-agent/internal/commands -> ../../..
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	installSh, err := os.ReadFile(filepath.Join(root, "install.sh"))
	if err != nil {
		t.Fatalf("read install.sh: %v", err)
	}

	for _, m := range refs {
		name := m[1]
		src := filepath.Join(root, "install", "systemd", name)
		if _, err := os.Stat(src); err != nil {
			t.Errorf("unit references /usr/local/libexec/dotnext/%s but infrastructure/systemd/%s does not exist (203/EXEC at runtime)", name, name)
		}
		if !strings.Contains(string(installSh), "infrastructure/systemd/"+name) {
			t.Errorf("install.sh never installs infrastructure/systemd/%s — generated unit would hit 203/EXEC on the host", name)
		}
	}
}
