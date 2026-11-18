package websites

import "testing"

func TestValidateAndDefaultDocRoot(t *testing.T) {
	t.Parallel()

	dr, err := DefaultDocRoot("alice", " Example.COM. ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dr != "/home/alice/web/example.com/public_html" {
		t.Fatalf("docroot=%q", dr)
	}

	for _, bad := range []string{"localhost", "-bad.com", "a..b.com", "bad;rm.com"} {
		if err := ValidateDomain(bad); err == nil {
			t.Errorf("expected ValidateDomain(%q) error", bad)
		}
	}
}
