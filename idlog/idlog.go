// Package idlog provides structured JSON logging with readable stack traces.
// Designed for Yandex Cloud Logging and standard-library-only environments.
package idlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// StreamName is the name of the log stream. Must be 1–63 characters.
// Set it once at application startup.
var StreamName = "app"

const maxStreamNameLen = 63

func init() {
	if len(StreamName) == 0 || len(StreamName) > maxStreamNameLen {
		panic("idlog: StreamName must be 1–63 characters")
	}
}

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

// Debug logs a DEBUG message only if enabled is true.
func Debug(message string, enabled bool) {
	if enabled {
		writeLog(LevelDebug, message)
	}
}

// Info logs an INFO message.
func Info(message string) {
	writeLog(LevelInfo, message)
}

// Warn logs a WARN message with a function-only stack trace.
func Warn(message string) {
	trace := getFormattedTrace(3) // skip Warn, writeLog and caller
	if trace != "" {
		message += " | TRACE: " + trace
	}
	writeLog(LevelWarn, message)
}

// Error logs an ERROR message without a function-only stack trace.
// Use AddTrace function to add a trace manually if needed.
func Error(err error) {
	if err == nil {
		err = fmt.Errorf("nil error passed to Error")
	}
	writeLog(LevelError, err.Error())
}

// Log logs based on the type of the argument:
//   - string → INFO
//   - error  → ERROR (without stack trace)
//   - other  → WARN (with stack trace)
func Log(v any) {
	switch val := v.(type) {
	case string:
		Info(val)
	case error:
		// Log error without stack trace
		writeLog(LevelError, val.Error())
	default:
		msg := fmt.Sprintf("unsupported log type: %T", v)
		// Log warn with stack trace
		trace := getFormattedTrace(3) // skip internal calls (Log, writeLog, and getFormattedTrace)
		if trace != "" {
			msg += " | TRACE: " + trace
		}
		writeLog(LevelWarn, msg)
	}
}

// AddTrace adds a human-readable function-only stack trace to an error.
func AddTrace(err error) error {
	if err == nil {
		return fmt.Errorf("nil error passed to AddTrace")
	}

	trace := getFormattedTrace(3) // skip AddTrace, internal calls and caller
	if trace != "" {
		return fmt.Errorf("%v | TRACE: %s", err, trace)
	}
	return err
}

// Fail logs a FATAL error with a function-only stack trace and exits the program.
func Fail(err error) {
	if err == nil {
		err = fmt.Errorf("nil error passed to Fail")
	}

	trace := getFormattedTrace(3) // skip Fail, writeLog and caller
	message := err.Error()
	if trace != "" {
		message += " | TRACE: " + trace
	}
	writeLog(LevelFatal, message)
	os.Exit(1)
}

// -- Internal helpers --

// getFormattedTrace returns a formatted function-only stack trace, skipping 'skip' top frames.
func getFormattedTrace(skip int) string {
	return formatFunctionOnlyTrace(getStackFrames(skip))
}

func writeLog(level, message string) {
	entry := struct {
		Message    string `json:"message"`
		Level      string `json:"level"`
		StreamName string `json:"stream_name"`
	}{
		Message:    message,
		Level:      level,
		StreamName: StreamName,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(entry) // adds \n automatically
	os.Stdout.Write(buf.Bytes())
}

// getStackFrames returns raw stack lines, skipping 'skip' top frames.
func getStackFrames(skip int) []string {
	buf := make([]byte, 1024*32)
	n := runtime.Stack(buf, false)
	lines := strings.Split(string(buf[:n]), "\n")

	i := 1 // skip "goroutine N [running]:"
	skipped := 0

	for i < len(lines) && skipped < skip {
		if strings.TrimSpace(lines[i]) == "" {
			i++
			continue
		}
		i++ // function line
		if i < len(lines) && strings.HasPrefix(lines[i], "\t") {
			i++ // file:line
		}
		skipped++
	}

	return lines[i:]
}

// formatFunctionOnlyTrace extracts only function names (without file paths or line numbers).
func formatFunctionOnlyTrace(stackLines []string) string {
	if len(stackLines) == 0 {
		return ""
	}

	var funcs []string
	for _, line := range stackLines {
		if strings.HasPrefix(line, "\t") {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Удаляем полный путь к пакету, оставляем только последнюю часть: pkg.Func → Func
		if idx := strings.LastIndex(trimmed, "/"); idx != -1 {
			trimmed = trimmed[idx+1:]
		}
		funcs = append(funcs, trimmed)
	}

	// Reverse to show natural call flow: caller >> callee
	for i, j := 0, len(funcs)-1; i < j; i, j = i+1, j-1 {
		funcs[i], funcs[j] = funcs[j], funcs[i]
	}

	// Удаляем дублирующиеся или слишком глубокие фреймы, если нужно
	// Например, оставим не более 5 уровней
	if len(funcs) > 5 {
		funcs = funcs[:5]
	}

	return strings.Join(funcs, " >> ")
}
