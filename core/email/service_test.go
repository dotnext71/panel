package mail

import (
	"errors"
	"testing"

	"github.com/dotnext-panel/dotnext/internal/mailaddr"
)

func TestValidateMailbox(t *testing.T) {
	t.Parallel()

	spec, err := ValidateMailbox("  Alice.Smith@Example.COM ", 2048, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.CanonicalAddr != "alice.smith@example.com" || spec.QuotaMB != 2048 {
		t.Fatalf("unexpected spec: %+v", spec)
	}

	if _, err := ValidateMailbox(mailaddr.DirectoryLocalPart+"@example.com", 500, true); !errors.Is(err, ErrReservedLocal) {
		t.Fatalf("expected ErrReservedLocal for DirectoryLocalPart, got %v", err)
	}
	if _, err := ValidateMailbox("postmaster@example.com", 500, false); !errors.Is(err, ErrReservedLocal) {
		t.Fatalf("expected ErrReservedLocal for postmaster on non-primary domain, got %v", err)
	}
	if _, err := ValidateMailbox("alice@example.com", -1, false); !errors.Is(err, ErrInvalidQuota) {
		t.Fatalf("expected ErrInvalidQuota, got %v", err)
	}
}

func TestNormalizeForwardTargets(t *testing.T) {
	t.Parallel()

	got, err := NormalizeForwardTargets([]string{"Bob@Example.org", "bob@example.org", "alice@example.net"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "bob@example.org" || got[1] != "alice@example.net" {
		t.Fatalf("unexpected normalized list: %v", got)
	}

	if _, err := NormalizeForwardTargets(nil); !errors.Is(err, ErrEmptyForwardList) {
		t.Fatalf("expected ErrEmptyForwardList, got %v", err)
	}
	if _, err := NormalizeForwardTargets([]string{"not-an-email"}); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("expected ErrInvalidAddress, got %v", err)
	}
}
