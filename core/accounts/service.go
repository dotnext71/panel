// Package accounts provides domain-level validation and transaction-safe
// account & hosting-package orchestration for dotNext.
package accounts

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"github.com/dotnext-panel/dotnext/internal/platform/database"
)

var (
	ErrInvalidUsername = errors.New("accounts: invalid linux username")
	ErrInvalidEmail    = errors.New("accounts: invalid email address")
	ErrReservedName    = errors.New("accounts: reserved system username")
)

var usernamePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,31}$`)

var reservedSystemUsers = map[string]struct{}{
	"root": {}, "daemon": {}, "bin": {}, "sys": {}, "sync": {},
	"www-data": {}, "mysql": {}, "postgres": {}, "redis": {},
	"dotnext": {}, "stalwart": {}, "pdns": {},
	"nobody": {}, "mail": {}, "postfix": {}, "dovecot": {},
}

// CreateAccountInput captures validated tenant creation attributes.
type CreateAccountInput struct {
	Username  string
	Email     string
	PackageID string
}

// Validate normalizes and checks tenant account inputs at the module boundary.
func (in *CreateAccountInput) Validate() error {
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.PackageID = strings.TrimSpace(in.PackageID)

	if !usernamePattern.MatchString(in.Username) {
		return fmt.Errorf("%w: must match ^[a-z][a-z0-9_-]{1,31}$", ErrInvalidUsername)
	}
	if _, reserved := reservedSystemUsers[in.Username]; reserved {
		return fmt.Errorf("%w: %q", ErrReservedName, in.Username)
	}
	if in.Email == "" || !strings.Contains(in.Email, "@") || strings.ContainsAny(in.Email, " \t\r\n") {
		return ErrInvalidEmail
	}
	return nil
}

// RunInTransaction executes multi-table account lifecycle work inside one atomic transaction.
func RunInTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return database.WithTx(ctx, db, fn)
}
