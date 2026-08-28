# idconfig

`idconfig` loads environment variables into a Go structure using field tags. It
uses only the standard library and does not log configuration values.

## Usage

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/idfactory/idtools/idconfig"
)

type Config struct {
	Environment string        `env:"APP_ENV" required:"true"`
	Name        string        `env:"APP_NAME" default:"idtools-example"`
	Debug       bool          `env:"DEBUG" default:"false"`
	Port        int           `env:"PORT" default:"8080"`
	Timeout     time.Duration `env:"TIMEOUT" default:"30s"`
}

func main() {
	config := &Config{}
	if err := idconfig.Load(config); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Debug: %t\n", config.Debug)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("Timeout: %s\n", config.Timeout)
}
```

Run it with the required value set:

```sh
APP_ENV=development go run ./examples/idconfig/
```

`Load` supports `string`, `bool`, `int`, `int64`, and `time.Duration`. Duration
values use Go syntax such as `30s`, `5m`, or `2h`.

- `env:"NAME"` selects the environment variable.
- `default:"value"` applies only when the variable is absent.
- `required:"true"` rejects a missing or empty variable.
- Fields without an `env` tag remain unchanged.

If several fields are invalid, `Load` returns all their errors together. It
does not partially update the target structure and does not include environment
values in errors.

A runnable example is available in
[`examples/idconfig`](../examples/idconfig/main.go).

## Optional `.env` file

`idconfig` intentionally reads only the process environment and therefore does
not add a `.env` parser or a third-party dependency to `idtools`. Applications
that need `.env` support can load the file before calling `idconfig.Load` with
[`godotenv`](https://github.com/joho/godotenv):

```sh
go get github.com/joho/godotenv
```

Add the following imports to the application configuration package:

```go
import (
	"errors"
	"fmt"
	"os"

	"github.com/idfactory/idtools/idconfig"
	"github.com/joho/godotenv"
)
```

Then load `.env` immediately before loading the configuration structure:

```go
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	config := &Config{}
	if err := idconfig.Load(config); err != nil {
		return nil, err
	}
	return config, nil
}
```

A missing `.env` file is ignored, while access and parsing errors are returned.
`godotenv.Load` preserves environment variables that are already set, so values
injected by the shell, a container, or the deployment platform take precedence
over local `.env` values. Keep `.env` in the consuming application's
`.gitignore`; commit a secret-free `.env.example` when configuration keys need
to be documented.
