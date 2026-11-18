package commands

import (
	"strings"
	"testing"
)

// GH #746: a systemd-run --property flag placed AFTER /usr/local/bin/dotnext-cli is
// handed to `dotnext migrate import` and dies with "unknown flag: --property",
// which broke every auto-create import (the only path that sets it). The
// EnvironmentFile property must sit before the command.
func TestImportSystemdRunArgs_PropertyBeforeCommand_GH746(t *testing.T) {
	args := importSystemdRunArgs(
		"01KYNJ4TMHTYQH133SKARBZ4R2", "demolab",
		[]string{"--property=EnvironmentFile=/run/dotnext-migrate-x.env"},
		[]string{"--target-email=a@b.test"},
	)

	idxProp, idxDotNext := -1, -1
	for i, a := range args {
		switch a {
		case "--property=EnvironmentFile=/run/dotnext-migrate-x.env":
			idxProp = i
		case "/usr/local/bin/dotnext-cli":
			idxDotNext = i
		}
	}
	if idxProp == -1 {
		t.Fatalf("--property not present in args: %v", args)
	}
	if idxDotNext == -1 {
		t.Fatalf("dotnext binary not present in args: %v", args)
	}
	if idxProp > idxDotNext {
		t.Errorf("--property (idx %d) must come BEFORE /usr/local/bin/dotnext-cli (idx %d); args=%v",
			idxProp, idxDotNext, args)
	}
}

func TestImportSystemdRunArgs_NoPropertyWhenNoPassword_GH746(t *testing.T) {
	args := importSystemdRunArgs("01KYNJ4TMHTYQH133SKARBZ4R2", "demolab", nil, nil)
	// GH746's actual invariant: every --property precedes the command —
	// one placed after /usr/local/bin/dotnext-cli lands in dotnext's argv and
	// dies with "unknown flag". (Since JAB-242 the baseline resource
	// props are ALWAYS present, so "zero properties" is no longer the
	// no-password shape.)
	cmdAt := -1
	for i, a := range args {
		if a == "/usr/local/bin/dotnext-cli" {
			cmdAt = i
			break
		}
	}
	if cmdAt < 0 {
		t.Fatalf("command not found in args=%v", args)
	}
	for _, a := range args[cmdAt:] {
		if strings.HasPrefix(a, "--property") {
			t.Errorf("--property after the command lands in dotnext argv (GH746); args=%v", args)
		}
	}
	// JAB-242: baseline resource controls are always present, before the command.
	joined := strings.Join(args[:cmdAt], " ")
	for _, want := range []string{"MemoryMax=", "TasksMax=", "CPUWeight=", "IOWeight="} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing resource prop %s before command; args=%v", want, args)
		}
	}
	// The command tail is still well-formed.
	joined = strings.Join(args, " ")
	if !strings.Contains(joined, "/usr/local/bin/dotnext-cli migrate import --job-id=01KYNJ4TMHTYQH133SKARBZ4R2 --target-user=demolab") {
		t.Errorf("malformed import command: %v", args)
	}
}
