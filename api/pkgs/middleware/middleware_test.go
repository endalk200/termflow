package middleware_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/endalk200/termflow-api/pkgs/middleware"
)

// TestLogger is a simple logger to capture logs in a string.
type TestLogger struct {
	logs []string
}

// Enabled is required by the slog.Handler interface. It controls whether
// logging should be performed at the given level and context.
func (tl *TestLogger) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle processes the log entry by appending it to the logs slice.
func (tl *TestLogger) Handle(_ context.Context, r slog.Record) error {
	var logEntry strings.Builder

	// Collect log fields into a string for verification
	logEntry.WriteString(fmt.Sprintf("level=%s ", r.Level.String()))
	r.Attrs(func(attr slog.Attr) bool {
		logEntry.WriteString(attr.Key)
		logEntry.WriteString("=")
		logEntry.WriteString(fmt.Sprintf("%v ", attr.Value))
		return true
	})

	// Append the log entry to the logs slice
	tl.logs = append(tl.logs, logEntry.String())
	return nil
}

// WithAttrs is required by the slog.Handler interface to return a new handler with the provided attributes.
func (tl *TestLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return tl
}

// WithGroup is required by the slog.Handler interface to return a new handler with the provided group name.
func (tl *TestLogger) WithGroup(name string) slog.Handler {
	return tl
}

func TestLoggingMiddleware(t *testing.T) {
	logger := &TestLogger{}                      // Create the test logger
	slogLogger := slog.New(logger)               // Create an slog.Logger using the test logger
	middleware := middleware.Logging(slogLogger) // Apply the middleware

	// Create a test HTTP handler that responds with "OK"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap the test handler with the logging middleware
	wrappedHandler := middleware(testHandler)

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	// Call the wrapped handler with the test request
	wrappedHandler.ServeHTTP(rr, req)

	// Assert that the response is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Ensure the logs contain the expected fields
	if len(logger.logs) == 0 {
		t.Fatal("expected logs to be written, but found none")
	}

	// Check that specific structured log fields exist in the logs
	logEntry := logger.logs[0]
	if !strings.Contains(logEntry, "method=GET") {
		t.Errorf("expected log to contain method=GET, got %s", logEntry)
	}
	if !strings.Contains(logEntry, "path=/test-path") {
		t.Errorf("expected log to contain path=/test-path, got %s", logEntry)
	}
	if !strings.Contains(logEntry, "status=200") {
		t.Errorf("expected log to contain status=200, got %s", logEntry)
	}
	if !strings.Contains(logEntry, "duration=") {
		t.Errorf("expected log to contain duration, got %s", logEntry)
	}
}
