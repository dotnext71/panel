package backups

import (
	"errors"
	"testing"
)

func TestValidateSchedulePolicy(t *testing.T) {
	t.Parallel()

	p, err := ValidateSchedulePolicy("DAILY", "01:00", "05:00", 14)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Cadence != "daily" || p.KeepCount != 14 {
		t.Fatalf("unexpected policy: %+v", p)
	}

	if _, err := ValidateSchedulePolicy("every-minute", "01:00", "05:00", 7); !errors.Is(err, ErrInvalidCadence) {
		t.Fatalf("expected ErrInvalidCadence, got %v", err)
	}
	if _, err := ValidateSchedulePolicy("daily", "25:99", "05:00", 7); !errors.Is(err, ErrInvalidWindow) {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
	if _, err := ValidateSchedulePolicy("daily", "01:00", "05:00", 0); !errors.Is(err, ErrInvalidRetention) {
		t.Fatalf("expected ErrInvalidRetention, got %v", err)
	}
}
