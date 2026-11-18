// Package jobs provides a supervised background worker manager for long-lived
// loops (reconcilers, audit consumer, metrics tickers, backup scheduler) with
// panic recovery, lifecycle status tracking, and graceful shutdown.
package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Runner is a long-lived background worker bound to ctx.
type Runner func(ctx context.Context) error

// WorkerState summarizes the live health of one supervised worker.
type WorkerState struct {
	Name      string     `json:"name"`
	Running   bool       `json:"running"`
	Restarts  int        `json:"restarts"`
	LastStart *time.Time `json:"last_start,omitempty"`
	LastError string     `json:"last_error,omitempty"`
}

type workerEntry struct {
	name         string
	fn           Runner
	restartDelay time.Duration
	state        WorkerState
}

// Supervisor manages named background workers under a single parent context.
type Supervisor struct {
	log     *slog.Logger
	mu      sync.RWMutex
	entries []*workerEntry
	wg      sync.WaitGroup
}

// NewSupervisor constructs a Supervisor. If log is nil, slog.Default() is used.
func NewSupervisor(log *slog.Logger) *Supervisor {
	if log == nil {
		log = slog.Default()
	}
	return &Supervisor{log: log}
}

// Register adds a named worker. restartDelay <= 0 means run once without
// automatic restart unless panicked/restarted; positive delays restart on
// non-cancellation exit.
func (s *Supervisor) Register(name string, restartDelay time.Duration, fn Runner) {
	if name == "" || fn == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, &workerEntry{
		name:         name,
		fn:           fn,
		restartDelay: restartDelay,
		state:        WorkerState{Name: name},
	})
}

// Start launches all registered workers and returns immediately. Call Wait()
// after cancelling ctx to drain cleanly.
func (s *Supervisor) Start(ctx context.Context) {
	s.mu.RLock()
	snapshot := make([]*workerEntry, len(s.entries))
	copy(snapshot, s.entries)
	s.mu.RUnlock()

	for _, entry := range snapshot {
		s.wg.Add(1)
		go s.runWorker(ctx, entry)
	}
}

// Wait blocks until every launched worker goroutine has exited.
func (s *Supervisor) Wait() {
	s.wg.Wait()
}

// Snapshot returns a copy of all worker states for monitoring endpoints.
func (s *Supervisor) Snapshot() []WorkerState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]WorkerState, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e.state)
	}
	return out
}

func (s *Supervisor) runWorker(ctx context.Context, e *workerEntry) {
	defer s.wg.Done()

	for {
		if ctx.Err() != nil {
			s.setStopped(e, "")
			return
		}

		now := time.Now().UTC()
		s.mu.Lock()
		e.state.Running = true
		e.state.LastStart = &now
		s.mu.Unlock()

		err := s.callSafely(ctx, e.fn)
		if ctx.Err() != nil {
			s.setStopped(e, "")
			return
		}

		errText := ""
		if err != nil {
			errText = err.Error()
			s.log.Warn("supervised worker exited with error", "worker", e.name, "err", err)
		}
		s.setStopped(e, errText)

		if e.restartDelay <= 0 {
			return
		}

		s.mu.Lock()
		e.state.Restarts++
		s.mu.Unlock()

		timer := time.NewTimer(e.restartDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (s *Supervisor) setStopped(e *workerEntry, lastErr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.state.Running = false
	e.state.LastError = lastErr
}

func (s *Supervisor) callSafely(ctx context.Context, fn Runner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return fn(ctx)
}
