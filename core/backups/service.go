// Package backups provides dotNext backup policy validation and manifest
// verification helpers built over internal/backup.
package backups

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dotnext-panel/dotnext/internal/backup"
)

var (
	ErrInvalidCadence   = errors.New("backups: invalid schedule cadence")
	ErrInvalidWindow    = errors.New("backups: invalid backup window HH:MM")
	ErrInvalidRetention = errors.New("backups: keep_count must be between 1 and 365")
)

// SchedulePolicy represents a validated backup schedule policy.
type SchedulePolicy struct {
	Cadence   string
	StartHHMM string
	EndHHMM   string
	KeepCount int
}

// ValidateSchedulePolicy verifies schedule cadence, optional window bounds, and retention count.
func ValidateSchedulePolicy(cadence, startHHMM, endHHMM string, keepCount int) (SchedulePolicy, error) {
	c := strings.ToLower(strings.TrimSpace(cadence))
	if !backup.ValidCadence(c) {
		return SchedulePolicy{}, fmt.Errorf("%w: %q", ErrInvalidCadence, cadence)
	}
	start := strings.TrimSpace(startHHMM)
	end := strings.TrimSpace(endHHMM)
	if start != "" && !backup.ValidHHMM(start) {
		return SchedulePolicy{}, fmt.Errorf("%w: start %q", ErrInvalidWindow, start)
	}
	if end != "" && !backup.ValidHHMM(end) {
		return SchedulePolicy{}, fmt.Errorf("%w: end %q", ErrInvalidWindow, end)
	}
	if keepCount < 1 || keepCount > 365 {
		return SchedulePolicy{}, ErrInvalidRetention
	}
	return SchedulePolicy{
		Cadence:   c,
		StartHHMM: start,
		EndHHMM:   end,
		KeepCount: keepCount,
	}, nil
}

// DecodeAccountManifest strictly decodes and validates a backup manifest payload.
func DecodeAccountManifest(raw []byte) (*backup.AccountManifest, error) {
	return backup.AccountManifestFromBytes(raw)
}
