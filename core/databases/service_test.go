package databases

import (
	"errors"
	"testing"
)

func TestNormalizeEngineAndValidateIdentifier(t *testing.T) {
	t.Parallel()

	if e, err := NormalizeEngine(""); err != nil || e != "mariadb" {
		t.Fatalf("got (%q, %v)", e, err)
	}
	if e, err := NormalizeEngine("postgresql"); err != nil || e != "postgres" {
		t.Fatalf("got (%q, %v)", e, err)
	}
	if _, err := NormalizeEngine("sqlite"); !errors.Is(err, ErrUnsupportedEngine) {
		t.Fatalf("expected ErrUnsupportedEngine, got %v", err)
	}

	if id, err := ValidateIdentifier("app_db_1"); err != nil || id != "app_db_1" {
		t.Fatalf("expected valid identifier, got (%q, %v)", id, err)
	}
	if _, err := ValidateIdentifier("mysql"); !errors.Is(err, ErrReservedIdentifier) {
		t.Fatalf("expected ErrReservedIdentifier for mysql, got %v", err)
	}
	if _, err := ValidateIdentifier("bad-name"); !errors.Is(err, ErrInvalidIdentifier) {
		t.Fatalf("expected ErrInvalidIdentifier, got %v", err)
	}
}
