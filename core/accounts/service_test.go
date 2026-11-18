package accounts

import (
	"errors"
	"testing"
)

func TestCreateAccountInputValidate(t *testing.T) {
	t.Parallel()

	ok := CreateAccountInput{Username: "  Alice_01 ", Email: " Alice@Example.com "}
	if err := ok.Validate(); err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	if ok.Username != "alice_01" || ok.Email != "alice@example.com" {
		t.Fatalf("unexpected normalization: %+v", ok)
	}

	for _, bad := range []CreateAccountInput{
		{Username: "root", Email: "a@b.com"},
		{Username: "dotnext", Email: "a@b.com"},
		{Username: "1bad", Email: "a@b.com"},
		{Username: "alice", Email: "not-an-email"},
	} {
		if err := bad.Validate(); err == nil {
			t.Errorf("expected error for %+v", bad)
		}
	}

	reserved := CreateAccountInput{Username: "dotnext", Email: "a@b.com"}
	if err := reserved.Validate(); !errors.Is(err, ErrReservedName) {
		t.Fatalf("expected ErrReservedName, got %v", err)
	}
}
