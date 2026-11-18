package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewJSONAndTextLevels(t *testing.T) {
	t.Parallel()

	var jsonBuf bytes.Buffer
	jlog := New("json", "warn", &jsonBuf)
	jlog.Info("ignored")
	jlog.Warn("captured", "component", "test")

	var decoded map[string]any
	if err := json.Unmarshal(jsonBuf.Bytes(), &decoded); err != nil {
		t.Fatalf("expected valid JSON log line, got %q: %v", jsonBuf.String(), err)
	}
	if decoded["msg"] != "captured" || decoded["level"] != "WARN" {
		t.Fatalf("unexpected decoded log: %v", decoded)
	}

	var textBuf bytes.Buffer
	tlog := New("text", "debug", &textBuf)
	tlog.Debug("trace-event")
	if !strings.Contains(textBuf.String(), "trace-event") {
		t.Fatalf("expected debug line in text logger, got %q", textBuf.String())
	}
}

func TestContextRoundTrip(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := New("json", "info", &buf).With("req_id", "01ABC")
	ctx := WithContext(context.Background(), l)
	FromContext(ctx).Info("ok")

	if !strings.Contains(buf.String(), "01ABC") {
		t.Fatalf("expected req_id in context logger output, got %q", buf.String())
	}
	if FromContext(context.Background()) != slog.Default() {
		t.Fatal("expected fallback to slog.Default() when ctx has no logger")
	}
}
