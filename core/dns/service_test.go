package dns

import (
	"errors"
	"testing"
)

func TestCheckCNAMEConflict(t *testing.T) {
	t.Parallel()

	existing := []RecordSummary{
		{ID: "1", Name: "www.example.com", Type: "A"},
		{ID: "2", Name: "alias.example.com", Type: "CNAME"},
	}

	if err := CheckCNAMEConflict("www.example.com.", "CNAME", "", existing); !errors.Is(err, ErrCNAMEConflict) {
		t.Fatalf("expected ErrCNAMEConflict adding CNAME over A, got %v", err)
	}
	if err := CheckCNAMEConflict("alias.example.com", "TXT", "", existing); !errors.Is(err, ErrCNAMEConflict) {
		t.Fatalf("expected ErrCNAMEConflict adding TXT over CNAME, got %v", err)
	}
	if err := CheckCNAMEConflict("alias.example.com", "CNAME", "2", existing); err != nil {
		t.Fatalf("expected self-update of CNAME to succeed, got %v", err)
	}
	if err := CheckCNAMEConflict("www.example.com", "AAAA", "", existing); err != nil {
		t.Fatalf("expected A+AAAA coexistence to succeed, got %v", err)
	}
}
