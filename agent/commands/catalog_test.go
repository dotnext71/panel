package commands

import (
	"strings"
	"testing"
)

func TestRegistryContract_NonEmptyAndNamespaced(t *testing.T) {
	t.Parallel()

	cmds := Default.Commands()
	if len(cmds) < 400 {
		t.Fatalf("expected >= 400 registered agent commands, got %d", len(cmds))
	}
	for _, cmd := range cmds {
		if strings.TrimSpace(cmd) != cmd || !strings.Contains(cmd, ".") {
			t.Errorf("invalid command name %q: must be trimmed and namespaced (<domain>.<verb>)", cmd)
		}
	}
}
