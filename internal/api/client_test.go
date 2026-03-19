package api

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogDebugPayloadWritesWhenDebugEnabled(t *testing.T) {
	oldDebug := os.Getenv("DEBUG")
	t.Cleanup(func() {
		if oldDebug == "" {
			os.Unsetenv("DEBUG")
			return
		}
		os.Setenv("DEBUG", oldDebug)
	})
	if err := os.Setenv("DEBUG", "1"); err != nil {
		t.Fatalf("Setenv returned error: %v", err)
	}

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe returned error: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = oldStderr
	})

	logDebugPayload([]byte(`{"code":"fmt.Println()"}`))

	if err := w.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Copy returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `DEBUG render payload: {"code":"fmt.Println()"}`) {
		t.Fatalf("unexpected debug output: %q", output)
	}
}

func TestLogDebugPayloadSkipsWhenDebugDisabled(t *testing.T) {
	oldDebug := os.Getenv("DEBUG")
	t.Cleanup(func() {
		if oldDebug == "" {
			os.Unsetenv("DEBUG")
			return
		}
		os.Setenv("DEBUG", oldDebug)
	})
	os.Unsetenv("DEBUG")

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe returned error: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = oldStderr
	})

	logDebugPayload([]byte(`{"code":"fmt.Println()"}`))

	if err := w.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Copy returned error: %v", err)
	}

	if buf.Len() != 0 {
		t.Fatalf("expected no debug output, got %q", buf.String())
	}
}
