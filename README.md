# idtools

`idtools` is a small Go module with reusable, standard-library-only utilities.
It currently contains two packages:

- [`idlog`](./docs/idlog.md) — a structured JSON logger with log levels and
  readable stack traces.
- [`idtoken`](./docs/idtoken.md) — cryptographically secure UUID and token
  generators built on `crypto/rand`.

The module requires Go 1.21 or newer and has no third-party runtime
dependencies.

## Installation

Add the module to an existing Go project:

```sh
go get github.com/idfactory/idtools@latest
```

Import only the package you need:

```go
import (
	"github.com/idfactory/idtools/idlog"
	"github.com/idfactory/idtools/idtoken"
)
```

See the package documentation linked above for API details and usage.

## Repository layout

```text
idlog/       Structured logging package and its tests
idtoken/     UUID and token package and its tests
examples/    Runnable example programs
docs/        Package documentation
.zed/        Optional Zed tasks and debugger configurations
```

## Running the examples

Clone the repository and run each example from the repository root:

```sh
go run ./examples/idtoken/
go run ./examples/idlog/
```

The source is in [`examples/idtoken`](./examples/idtoken/main.go) and
[`examples/idlog`](./examples/idlog/main.go). You can copy the relevant calls
into your own program and adapt their arguments and error handling.

The `idlog` example intentionally ends with `idlog.Fail`. It prints a `FATAL`
log entry and exits with status 1; that non-zero exit is expected for this
example.

## Tests and checks

Run the complete test suite from the repository root:

```sh
go test ./...
```

Useful additional checks are:

```sh
go test -v ./...      # verbose test output
go test -cover ./...  # package-level test coverage
go vet ./...          # static analysis
```

Run all tests and `go vet ./...` before opening a pull request.

## Using Zed

The repository includes optional project configuration for
[Zed](https://zed.dev/) in [`.zed/`](./.zed/):

- [`.zed/tasks.json`](./.zed/tasks.json) provides tasks for both examples,
  tests, coverage, and `go vet`.
- [`.zed/debug.json`](./.zed/debug.json) provides Delve launch configurations
  for both examples.
- [`.zed/settings.json`](./.zed/settings.json) contains project-local editor
  settings.

Open the repository root in Zed. To run a task, open the command palette
(`Cmd+Shift+P` on macOS or `Ctrl+Shift+P` on Linux), choose `task: spawn`, then
select a task such as `Run Tests` or `Run Example - idtoken`. To debug an
example, run `debugger: start` from the command palette and select one of the
`Debug example - ...` launch configurations. Zed has built-in Go debugging
support through Delve.

Zed is not required. With any other editor—or no editor-specific
integration—run the commands in [Running the examples](#running-the-examples)
and [Tests and checks](#tests-and-checks) from a terminal. The `.zed/` files do
not affect builds, examples, or tests run by the Go toolchain.

## Contributing

The contribution guide is [`CONTRIBUTING.md`](./CONTRIBUTING.md) in the
repository root. Before making a change:

1. Read the guide and check the existing GitHub issues.
2. Open an issue to discuss the proposed change with the maintainer.
3. After the approach is agreed, create a branch, make the change, and run the
   tests and checks above.
4. Open a pull request to `main` with the reason for the change, verification
   steps, and links to related issues.

Use the guide as the source of truth for the full contribution workflow,
review expectations, and release instructions.

## License

This project is available under the [MIT License](./LICENSE).
