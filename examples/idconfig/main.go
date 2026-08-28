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
