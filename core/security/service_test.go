package security

import (
	"errors"
	"testing"
)

func TestParseIPOrCIDR(t *testing.T) {
	t.Parallel()

	p, err := ParseIPOrCIDR("203.0.113.9")
	if err != nil || p.String() != "203.0.113.9/32" {
		t.Fatalf("unexpected prefix: %v, err=%v", p, err)
	}
	p2, err := ParseIPOrCIDR("2001:db8::1/64")
	if err != nil || p2.String() != "2001:db8::/64" {
		t.Fatalf("unexpected masked IPv6 prefix: %v, err=%v", p2, err)
	}
	if _, err := ParseIPOrCIDR("999.999.1.1"); !errors.Is(err, ErrInvalidCIDROrIP) {
		t.Fatalf("expected ErrInvalidCIDROrIP, got %v", err)
	}
}

func TestValidateSSHPublicKey(t *testing.T) {
	t.Parallel()

	k, err := ValidateSSHPublicKey("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMockKeyData alice@host")
	if err != nil || k == "" {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := ValidateSSHPublicKey("ssh-dss AAAAB3Nza..."); !errors.Is(err, ErrInvalidSSHKey) {
		t.Fatalf("expected ErrInvalidSSHKey for ssh-dss, got %v", err)
	}
	if _, err := ValidateSSHPublicKey("ssh-ed25519 AAA\ncommand=\"rm -rf\""); !errors.Is(err, ErrInvalidSSHKey) {
		t.Fatalf("expected ErrInvalidSSHKey for newline injection, got %v", err)
	}
}
