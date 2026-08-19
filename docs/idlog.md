# idlog — Structured JSON Logger for Go

<img src="./src/idlog_logo.png" alt="idlog logo" width="250">

A minimal, zero-dependency logger for **structured JSON** output with **readable stack traces**.  
Built with **only the standard library**, designed for **Yandex Cloud Logging** and personal projects.

## ✅ Features

- **Single-line JSON** output
- Log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`
- Global `stream_name` (1–63 chars) for logical grouping
- Type-driven logging via `Log(any)`:
  - `string` → `INFO`
  - `error` → `ERROR` (no trace by default)
  - any other type → `WARN` + **embedded function trace**
- `Debug(message string, enabled bool)` — respects debug flag
- `Fail(error)` → logs `FATAL`, emits `TRACE`, exits program
- `AddTrace(error)` → adds **human-readable function trace** to an error manually
- **Only standard library used**: `encoding/json`, `os`, `runtime`, `strings`, `fmt`, `bytes`
- **Human-readable trace** embedded directly in the `message` field, showing **caller >> callee** flow

## 🚀 Usage

### Install

```bash
go get github.com/idfactory/idtools@latest
```

### Use in code

```go
package main

import (
	"fmt"

	"github.com/idfactory/idtools/idlog"
)

func main() {
	idlog.StreamName = "my-service"

	idlog.Log("Start example")                   //INFO example
	idlog.Log(fmt.Errorf("error"))               //ERROR example (no trace)
	idlog.Log(1)                                 //WARN example (with trace)
	idlog.Debug("debug example", true)           //DEBUG example

	// Manually add trace to errors
	err := fmt.Errorf("some error")
	errWithTrace := idlog.AddTrace(err)
	idlog.Log(errWithTrace)                      //ERROR with manual trace

	idlog.Fail(fmt.Errorf("critical error"))		 //FAIL example
}
```

## 📤 Output Examples

### ERROR + TRACE

```json
{"message":"some error | TRACE: main.main()","level":"ERROR","stream_name":"my-service"}
```

### DEBUG (enabled)
```json
{"message":"debug example","level":"DEBUG","stream_name":"my-service"}
```

## ☁️ Yandex Cloud Logging

- Logs to stdout → automatically captured by YCL.
- Filter by level, `stream_name`, or `message` in Yandex UI/CLI.
- No extra setup needed.

## 📦 Module Info

- No dependencies
- Go 1.26.5+
- License: MIT
