package idlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// logEntry represents a single log entry for testing purposes
type logEntry struct {
	Message    string `json:"message"`
	Level      string `json:"level"`
	StreamName string `json:"stream_name"`
}

// captureStdout captures output written to os.Stdout during function execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// parseLogEntry parses a single JSON log entry from the output
func parseLogEntry(output string) (*logEntry, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, fmt.Errorf("empty output")
	}

	var entry logEntry
	err := json.Unmarshal([]byte(output), &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// parseMultipleLogEntries parses multiple JSON log entries from the output
func parseMultipleLogEntries(output string) ([]logEntry, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var entries []logEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry logEntry
		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func TestLogWithString(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Log("test message")
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	expected := logEntry{
		Message:    "test message",
		Level:      LevelInfo,
		StreamName: "test-stream",
	}

	if entry.Message != expected.Message {
		t.Errorf("Expected message '%s', got '%s'", expected.Message, entry.Message)
	}
	if entry.Level != expected.Level {
		t.Errorf("Expected level '%s', got '%s'", expected.Level, entry.Level)
	}
	if entry.StreamName != expected.StreamName {
		t.Errorf("Expected stream name '%s', got '%s'", expected.StreamName, entry.StreamName)
	}
}

func TestLogWithError(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	testErr := fmt.Errorf("test error")

	output := captureStdout(func() {
		Log(testErr)
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	expected := logEntry{
		Message:    "test error",
		Level:      LevelError,
		StreamName: "test-stream",
	}

	if entry.Message != expected.Message {
		t.Errorf("Expected message '%s', got '%s'", expected.Message, entry.Message)
	}
	if entry.Level != expected.Level {
		t.Errorf("Expected level '%s', got '%s'", expected.Level, entry.Level)
	}
	if entry.StreamName != expected.StreamName {
		t.Errorf("Expected stream name '%s', got '%s'", expected.StreamName, entry.StreamName)
	}

	// Check that no trace is included for errors
	if strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Error should not include trace by default, but got: %s", entry.Message)
	}
}

func TestLogWithOtherType(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Log(42) // Using an integer to trigger the 'other' case
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Level != LevelWarn {
		t.Errorf("Expected level '%s', got '%s'", LevelWarn, entry.Level)
	}
	if entry.StreamName != "test-stream" {
		t.Errorf("Expected stream name 'test-stream', got '%s'", entry.StreamName)
	}

	// Check that trace is included for non-string, non-error types
	if !strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Non-string/non-error type should include trace, but got: %s", entry.Message)
	}

	// Check that the message indicates unsupported type
	if !strings.Contains(entry.Message, "unsupported log type: int") {
		t.Errorf("Message should indicate unsupported type, but got: %s", entry.Message)
	}
}

func TestDebug(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Debug("debug message", true)
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	expected := logEntry{
		Message:    "debug message",
		Level:      LevelDebug,
		StreamName: "test-stream",
	}

	if entry.Message != expected.Message {
		t.Errorf("Expected message '%s', got '%s'", expected.Message, entry.Message)
	}
	if entry.Level != expected.Level {
		t.Errorf("Expected level '%s', got '%s'", expected.Level, entry.Level)
	}
	if entry.StreamName != expected.StreamName {
		t.Errorf("Expected stream name '%s', got '%s'", expected.StreamName, entry.StreamName)
	}

	// Check that no trace is included for debug
	if strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Debug should not include trace, but got: %s", entry.Message)
	}
}

func TestDebugDisabled(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Debug("debug message", false) // disabled
	})

	// When debug is disabled, nothing should be logged
	if output != "" {
		t.Errorf("Expected no output when debug is disabled, but got: %s", output)
	}
}

func TestInfo(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Info("info message")
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	expected := logEntry{
		Message:    "info message",
		Level:      LevelInfo,
		StreamName: "test-stream",
	}

	if entry.Message != expected.Message {
		t.Errorf("Expected message '%s', got '%s'", expected.Message, entry.Message)
	}
	if entry.Level != expected.Level {
		t.Errorf("Expected level '%s', got '%s'", expected.Level, entry.Level)
	}
	if entry.StreamName != expected.StreamName {
		t.Errorf("Expected stream name '%s', got '%s'", expected.StreamName, entry.StreamName)
	}

	// Check that no trace is included for info
	if strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Info should not include trace, but got: %s", entry.Message)
	}
}

func TestWarn(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Warn("warning message")
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Level != LevelWarn {
		t.Errorf("Expected level '%s', got '%s'", LevelWarn, entry.Level)
	}
	if entry.StreamName != "test-stream" {
		t.Errorf("Expected stream name 'test-stream', got '%s'", entry.StreamName)
	}

	// Check that trace is included for warn
	if !strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Warning should include trace, but got: %s", entry.Message)
	}
}

func TestError(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	testErr := fmt.Errorf("test error")

	output := captureStdout(func() {
		Error(testErr)
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	expected := logEntry{
		Message:    "test error",
		Level:      LevelError,
		StreamName: "test-stream",
	}

	if entry.Message != expected.Message {
		t.Errorf("Expected message '%s', got '%s'", expected.Message, entry.Message)
	}
	if entry.Level != expected.Level {
		t.Errorf("Expected level '%s', got '%s'", expected.Level, entry.Level)
	}
	if entry.StreamName != expected.StreamName {
		t.Errorf("Expected stream name '%s', got '%s'", expected.StreamName, entry.StreamName)
	}

	// Check that no trace is included for error by default
	if strings.Contains(entry.Message, "TRACE:") {
		t.Errorf("Error should not include trace by default, but got: %s", entry.Message)
	}
}

func TestAddTrace(t *testing.T) {
	testErr := fmt.Errorf("test error")

	resultErr := AddTrace(testErr)
	resultMsg := resultErr.Error()

	// Check that the trace was added
	if !strings.Contains(resultMsg, "TRACE:") {
		t.Errorf("AddTrace should add trace to error, but got: %s", resultMsg)
	}

	// Check that the original error message is preserved
	if !strings.Contains(resultMsg, "test error") {
		t.Errorf("AddTrace should preserve original error message, but got: %s", resultMsg)
	}
}

func TestAddTraceWithNil(t *testing.T) {
	resultErr := AddTrace(nil)
	resultMsg := resultErr.Error()

	// Check that a specific error message is returned for nil
	if !strings.Contains(resultMsg, "nil error passed to AddTrace") {
		t.Errorf("AddTrace with nil should return specific error, but got: %s", resultMsg)
	}
}

func TestErrorWithNil(t *testing.T) {
	originalStreamName := StreamName
	StreamName = "test-stream"
	defer func() {
		StreamName = originalStreamName
	}()

	output := captureStdout(func() {
		Error(nil) // This should handle nil gracefully
	})

	entry, err := parseLogEntry(output)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	// Check that a specific error message is logged when Error is called with nil
	if !strings.Contains(entry.Message, "nil error passed to Error") {
		t.Errorf("Error with nil should log specific message, but got: %s", entry.Message)
	}
	if entry.Level != LevelError {
		t.Errorf("Expected level '%s', got '%s'", LevelError, entry.Level)
	}
	if entry.StreamName != "test-stream" {
		t.Errorf("Expected stream name 'test-stream', got '%s'", entry.StreamName)
	}
}

func TestFormatFunctionOnlyTrace(t *testing.T) {
	stackLines := []string{
		"github.com/idfactory/idtools/idlog.Warn(...) ",
		"\t/project/idlog/idlog.go:42 +0x10",
		"main.main()",
		"\t/project/main.go:10 +0x20",
	}

	got := formatFunctionOnlyTrace(stackLines)
	want := "main.main() >> idlog.Warn(...)"
	if got != want {
		t.Fatalf("Expected trace %q, got %q", want, got)
	}
}
