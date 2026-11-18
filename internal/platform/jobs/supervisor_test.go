package jobs

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestSupervisorPanicRecoveryAndGracefulStop(t *testing.T) {
	t.Parallel()

	sup := NewSupervisor(nil)
	var runs atomic.Int32
	sup.Register("flaky", 10*time.Millisecond, func(ctx context.Context) error {
		n := runs.Add(1)
		if n == 1 {
			panic("boom")
		}
		if n == 2 {
			return errors.New("transient")
		}
		<-ctx.Done()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	sup.Start(ctx)

	deadline := time.Now().Add(2 * time.Second)
	for runs.Load() < 3 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if runs.Load() < 3 {
		t.Fatalf("expected worker to recover from panic and error, runs=%d", runs.Load())
	}

	snap := sup.Snapshot()
	if len(snap) != 1 || snap[0].Restarts < 2 || !snap[0].Running {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}

	cancel()
	sup.Wait()

	snapAfter := sup.Snapshot()
	if snapAfter[0].Running {
		t.Fatalf("expected worker stopped after cancel, got %+v", snapAfter[0])
	}
}
