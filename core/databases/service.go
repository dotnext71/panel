// Package databases provides engine validation and safe identifier normalization
// across MariaDB and PostgreSQL tenant databases.
package databases

import (
	"errors"
	"regexp"
	"strings"

	"github.com/dotnext-panel/dotnext/internal/dbreserve"
)

var (
	ErrUnsupportedEngine  = errors.New("databases: unsupported engine")
	ErrInvalidIdentifier  = errors.New("databases: invalid database or user identifier")
	ErrReservedIdentifier = errors.New("databases: reserved system identifier")
)

var identPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,62}$`)

// NormalizeEngine validates "mariadb" (default) or "postgres".
func NormalizeEngine(engine string) (string, error) {
	e := strings.ToLower(strings.TrimSpace(engine))
	if e == "" || e == "mysql" {
		return "mariadb", nil
	}
	if e == "postgresql" {
		e = "postgres"
	}
	if e != "mariadb" && e != "postgres" {
		return "", ErrUnsupportedEngine
	}
	return e, nil
}

// ValidateIdentifier verifies a safe SQL identifier and rejects reserved engine names.
func ValidateIdentifier(ident string) (string, error) {
	id := strings.ToLower(strings.TrimSpace(ident))
	if !identPattern.MatchString(id) {
		return "", ErrInvalidIdentifier
	}
	if dbreserve.Database(id) || dbreserve.User(id) {
		return "", ErrReservedIdentifier
	}
	return id, nil
}
