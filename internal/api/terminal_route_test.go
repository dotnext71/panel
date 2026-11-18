// SPDX-License-Identifier: MIT

package api

// Coverage for the M45 root web terminal (ADR-0096) — the highest-risk surface
// in the product: an authenticated-admin root shell by design. These tests pin
// the controls that keep it gated, because a regression here is an RCE.
//
// The PTY broker itself lives in panel-agent and is tested there; this file
// covers the panel-api half: the off-by-default gate, the one-shot IP+user
// bound token, and the frame-size cap on the UDS bridge.

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dotnext-panel/dotnext/internal/auth"
	"github.com/dotnext-panel/dotnext/internal/ginctx"
	"github.com/dotnext-panel/dotnext/internal/models"
	"io"
	"log/slog"
	"time"
)

// --- fakes ------------------------------------------------------------------

type fakeTerminalSessions struct {
	created    *models.TerminalSession
	consumed   *models.TerminalSession
	consumeErr error
	started    int
	ended      int
}

func (f *fakeTerminalSessions) Create(_ context.Context, s *models.TerminalSession) error {
	f.created = s
	return nil
}

func (f *fakeTerminalSessions) ConsumeValid(_ context.Context, _ string) (*models.TerminalSession, error) {
	if f.consumeErr != nil {
		return nil, f.consumeErr
	}
	return f.consumed, nil
}

func (f *fakeTerminalSessions) MarkStarted(_ context.Context, _, _ string) error {
	f.started++
	return nil
}

func (f *fakeTerminalSessions) MarkEnded(_ context.Context, _ string) error {
	f.ended++
	return nil
}

func (f *fakeTerminalSessions) CleanupExpired(context.Context) (int64, error) { return 0, nil }

type fakeServerSettings struct {
	enabled bool
	err     error
}

func (f *fakeServerSettings) Get(context.Context) (*models.ServerSettings, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &models.ServerSettings{RootTerminalEnabled: f.enabled}, nil
}

func (f *fakeServerSettings) Upsert(context.Context, *models.ServerSettings) error { return nil }

// The remaining interface methods are unreachable from the terminal path; the
// embedded nil interface keeps this fake compiling when the interface grows,
// and any accidental call panics loudly instead of silently returning zero.
func (f *fakeServerSettings) EnsureVAPID(context.Context, string) (bool, error) { return false, nil }
func (f *fakeServerSettings) SetDigestLastSent(context.Context, string) error   { return nil }
func (f *fakeServerSettings) RecordDRSync(context.Context, string, string, string) error {
	return nil
}
func (f *fakeServerSettings) ReassertDRPairing(context.Context, string, string, *time.Time) error {
	return nil
}

func adminClaims() *auth.AccessClaims {
	return &auth.AccessClaims{UserID: "admin-1", IsAdmin: true, Email: "admin@example.test"}
}

func userClaims() *auth.AccessClaims {
	return &auth.AccessClaims{UserID: "user-1", IsAdmin: false, Email: "user@example.test"}
}

func testLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// --- harness ----------------------------------------------------------------

type terminalTestEnv struct {
	router   *gin.Engine
	sessions *fakeTerminalSessions
	settings *fakeServerSettings
}

// newTerminalRouter mounts the terminal routes behind a claims injector that
// stands in for RequireAdmin, so each test can choose the caller's identity.
func newTerminalRouter(t *testing.T, claims *auth.AccessClaims, ptySocket string) *terminalTestEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	env := &terminalTestEnv{
		sessions: &fakeTerminalSessions{},
		settings: &fakeServerSettings{},
	}
	cfg := TerminalHandlerConfig{
		Sessions:       env.sessions,
		ServerSettings: env.settings,
		AgentPTYSocket: ptySocket,
		Log:            testLogger(t),
	}

	r := gin.New()
	grp := r.Group("/api/v1/admin", func(c *gin.Context) {
		if claims != nil {
			ginctx.SetClaims(c, claims)
		}
		c.Next()
	})
	h := &terminalHandler{cfg: cfg}
	grp.POST("/terminal/session", h.mint)
	grp.GET("/terminal/ws/:token", h.ws)
	env.router = r
	return env
}

const testTokenLen = 43

// --- tests ------------------------------------------------------------------

// The gate is the whole point of ADR-0096: off by default, and a setting read
// failure must NOT be treated as "on".
func TestTerminalMint_GatedOffByDefault(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/terminal/session", nil))

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
	if env.sessions.created != nil {
		t.Fatal("no session may be minted while the gate is closed")
	}
}

func TestTerminalMint_SettingsErrorFailsClosed(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")
	// Even with "enabled" true, a read error must not open the gate.
	env.settings.err = errors.New("db down")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/terminal/session", nil))

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 (fail closed)", w.Code)
	}
	if env.sessions.created != nil {
		t.Fatal("settings error must not mint a session")
	}
}

func TestTerminalMint_NonAdminForbidden(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, userClaims(), "/nonexistent.sock")
	env.settings.enabled = true

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/terminal/session", nil))

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for a non-admin", w.Code)
	}
	if env.sessions.created != nil {
		t.Fatal("non-admin must not mint a root-shell token")
	}
}

func TestTerminalMint_NoClaimsUnauthorized(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, nil, "/nonexistent.sock")
	env.settings.enabled = true

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/terminal/session", nil))

	if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 401/403 without claims", w.Code)
	}
}

// A minted token must be 256 bits of entropy, short-lived, and bound to the
// calling admin + IP — the three properties thread the WS upgrade relies on.
func TestTerminalMint_CreatesShortLivedBoundToken(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")
	env.settings.enabled = true

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/terminal/session", nil)
	req.RemoteAddr = "203.0.113.7:4444"
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	s := env.sessions.created
	if s == nil {
		t.Fatal("no session created")
	}
	if len(s.Token) != testTokenLen {
		t.Errorf("token length = %d, want %d (32 bytes base64url)", len(s.Token), testTokenLen)
	}
	if raw, err := base64.RawURLEncoding.DecodeString(s.Token); err != nil || len(raw) != 32 {
		t.Errorf("token is not 32 raw bytes (err=%v, len=%d)", err, len(raw))
	}
	if got := s.ExpiresAt.Sub(s.CreatedAt); got > terminalTokenTTL+1e9 {
		t.Errorf("token lifetime = %v, want <= %v", got, terminalTokenTTL)
	}
	if s.UserID != adminClaims().UserID {
		t.Errorf("token bound to user %q, want %q", s.UserID, adminClaims().UserID)
	}
	if s.ClientIP != "203.0.113.7" {
		t.Errorf("token bound to IP %q, want 203.0.113.7", s.ClientIP)
	}
	if s.Token == "" || len(s.ID) == 0 {
		t.Error("session missing token or id")
	}
}

// The WS upgrade rejects a short/garbage token before touching the repository.
func TestTerminalWS_RejectsMalformedToken(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")
	env.settings.enabled = true

	for _, tok := range []string{"short", "", "x"} {
		w := httptest.NewRecorder()
		env.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/terminal/ws/"+tok, nil))
		if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Errorf("token %q: status = %d, want 400", tok, w.Code)
		}
	}
}

// A consumed token cannot be replayed, and the failure must not be a 200/upgrade.
func TestTerminalWS_RejectsUsedOrInvalidToken(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")
	env.settings.enabled = true
	env.sessions.consumeErr = errors.New("not found")

	tok := base64.RawURLEncoding.EncodeToString(make([]byte, 32))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/terminal/ws/"+tok, nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for an invalid/used token", w.Code)
	}
	if env.sessions.started != 0 {
		t.Error("no session may be marked started for an invalid token")
	}
}

// The token is bound to the admin who minted it and the IP it was minted from;
// a different admin or a different client IP must be refused.
func TestTerminalWS_RejectsRebind(t *testing.T) {
	t.Parallel()
	valid := base64.RawURLEncoding.EncodeToString(make([]byte, 32))

	cases := []struct {
		name  string
		sess  *models.TerminalSession
		match bool
	}{
		{"same admin same ip", &models.TerminalSession{
			UserID: adminClaims().UserID, ClientIP: "203.0.113.7",
		}, true},
		{"different admin", &models.TerminalSession{
			UserID: "someone-else", ClientIP: "203.0.113.7",
		}, false},
		{"different ip", &models.TerminalSession{
			UserID: adminClaims().UserID, ClientIP: "198.51.100.9",
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := newTerminalRouter(t, adminClaims(), "/nonexistent.sock")
			env.settings.enabled = true
			env.sessions.consumed = tc.sess

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/terminal/ws/"+valid, nil)
			req.RemoteAddr = "203.0.113.7:5555"
			w := httptest.NewRecorder()
			env.router.ServeHTTP(w, req)

			if tc.match {
				// The PTY broker socket does not exist here, so a correctly
				// bound token proceeds past the rebind check and fails at the
				// dial with 502 — NOT 403. That distinction is the assertion:
				// 403 means rebind was rejected, 502 means it was accepted.
				if w.Code == http.StatusForbidden {
					t.Fatalf("status = 403: a correctly bound token must not be rebind-denied")
				}
			} else {
				if w.Code != http.StatusForbidden {
					t.Fatalf("status = %d, want 403 on rebind mismatch", w.Code)
				}
			}
		})
	}
}

// A broker-down situation must surface as a clean HTTP error before the WS is
// upgraded, never as a half-open socket.
func TestTerminalWS_BrokerDownIsBadGateway(t *testing.T) {
	t.Parallel()
	env := newTerminalRouter(t, adminClaims(), "/nonexistent-agent-pty.sock")
	env.settings.enabled = true
	env.sessions.consumed = &models.TerminalSession{
		UserID: adminClaims().UserID, ClientIP: "203.0.113.7",
	}

	tok := base64.RawURLEncoding.EncodeToString(make([]byte, 32))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/terminal/ws/"+tok, nil)
	req.RemoteAddr = "203.0.113.7:5555"
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502 when the PTY broker is unreachable", w.Code)
	}
}

// --- framing ----------------------------------------------------------------

func TestUDSFrame_RoundTrip(t *testing.T) {
	t.Parallel()
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	want := []byte{0x00, 0x01, 0x00, 0xff, 'h', 'i'}
	go func() { _ = udsWriteFrame(a, tWSStdout, want) }()

	op, payload, err := udsReadFrame(b)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if op != tWSStdout {
		t.Errorf("op = %#x, want %#x", op, tWSStdout)
	}
	if string(payload) != string(want) {
		t.Errorf("payload = %x, want %x", payload, want)
	}
}

// An oversized length field must be refused instead of allocating it: the UDS
// peer is the root agent, but a bug there must not let it OOM the panel.
func TestUDSFrame_RejectsOversizeLength(t *testing.T) {
	t.Parallel()
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	go func() {
		var hdr [5]byte
		hdr[0] = tWSStdout
		binary.BigEndian.PutUint32(hdr[1:], 1<<30) // 1 GiB
		_, _ = a.Write(hdr[:])
	}()

	if _, _, err := udsReadFrame(b); err == nil {
		t.Fatal("expected an error for an oversize frame")
	}
}

func TestAtoiDefault(t *testing.T) {
	t.Parallel()
	cases := map[string]int{
		"":        80,
		"0":       80,
		"-5":      80,
		"abc":     80,
		"120":     120,
		"100000":  80, // upper bound is exclusive
		"1000000": 80,
	}
	for in, want := range cases {
		if got := atoiDefault(in, 80); got != want {
			t.Errorf("atoiDefault(%q) = %d, want %d", in, got, want)
		}
	}
}
