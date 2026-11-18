// Package mail provides dotNext mail domain validation and boundary normalization,
// wrapping internal/mailaddr and enforcing quota/routing constraints.
package mail

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dotnext-panel/dotnext/internal/mailaddr"
)

var (
	ErrInvalidAddress   = errors.New("mail: invalid email address")
	ErrReservedLocal    = errors.New("mail: reserved mailbox local-part")
	ErrInvalidQuota     = errors.New("mail: quota_mb must be >= 0")
	ErrTooManyTargets   = errors.New("mail: forwarder exceeds max target count")
	ErrEmptyForwardList = errors.New("mail: forwarder must target at least one address")
)

const maxForwardTargets = 50

// MailboxSpec represents validated mailbox creation inputs.
type MailboxSpec struct {
	Local         string
	Domain        string
	CanonicalAddr string
	QuotaMB       int64
}

// ValidateMailbox canonicalizes address and checks reserved local-parts and quota bounds.
func ValidateMailbox(rawAddr string, quotaMB int64, isPrimaryPanelDomain bool) (MailboxSpec, error) {
	if quotaMB < 0 {
		return MailboxSpec{}, fmt.Errorf("%w: %d", ErrInvalidQuota, quotaMB)
	}
	local, domain, err := mailaddr.Canonicalise(strings.TrimSpace(rawAddr))
	if err != nil {
		return MailboxSpec{}, fmt.Errorf("%w: %v", ErrInvalidAddress, err)
	}
	if err := mailaddr.CheckNotReservedOn(local, isPrimaryPanelDomain); err != nil {
		return MailboxSpec{}, fmt.Errorf("%w: %v", ErrReservedLocal, err)
	}
	return MailboxSpec{
		Local:         local,
		Domain:        domain,
		CanonicalAddr: local + "@" + domain,
		QuotaMB:       quotaMB,
	}, nil
}

// NormalizeForwardTargets canonicalizes and deduplicates external forwarder recipients.
func NormalizeForwardTargets(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyForwardList
	}
	if len(raw) > maxForwardTargets {
		return nil, fmt.Errorf("%w: max %d", ErrTooManyTargets, maxForwardTargets)
	}
	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))
	for _, entry := range raw {
		canon, err := mailaddr.CanonicaliseExternal(strings.TrimSpace(entry))
		if err != nil {
			return nil, fmt.Errorf("%w: %q: %v", ErrInvalidAddress, entry, err)
		}
		if _, dup := seen[canon]; dup {
			continue
		}
		seen[canon] = struct{}{}
		out = append(out, canon)
	}
	return out, nil
}
