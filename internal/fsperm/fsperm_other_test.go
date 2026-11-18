//go:build !linux

package fsperm

import (
	"errors"
	"testing"
)

func TestNonLinuxStubsFailClosed(t *testing.T) {
	t.Parallel()

	if err := GroupSetgidTree("/home/u", 33); !errors.Is(err, errLinuxOnly) {
		t.Fatalf("expected errLinuxOnly, got %v", err)
	}
	if err := RepairDocrootGroup("/home/u", "/home/u/web", 33); !errors.Is(err, errLinuxOnly) {
		t.Fatalf("expected errLinuxOnly, got %v", err)
	}
	if _, err := ReownTree("/home/u", "/home/u/web", nil, 1000, nil, nil); !errors.Is(err, errLinuxOnly) {
		t.Fatalf("expected errLinuxOnly, got %v", err)
	}
}
