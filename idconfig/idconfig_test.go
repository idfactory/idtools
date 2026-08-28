package idconfig

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("APP_DEBUG", "true")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_TIMEOUT", "45s")
	t.Setenv("APP_NAME", "orders")

	type config struct {
		Debug    bool          `env:"APP_DEBUG" default:"false"`
		Port     int           `env:"APP_PORT" default:"80"`
		Timeout  time.Duration `env:"APP_TIMEOUT" default:"5s"`
		Name     string        `env:"APP_NAME" required:"true"`
		Workers  int64         `env:"APP_WORKERS" default:"4"`
		Existing string
	}

	got := config{Existing: "keep"}
	if err := Load(&got); err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	want := config{
		Debug:    true,
		Port:     8080,
		Timeout:  45 * time.Second,
		Name:     "orders",
		Workers:  4,
		Existing: "keep",
	}
	if got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{
		"IDCONFIG_TEST_DEFAULT_STRING",
		"IDCONFIG_TEST_DEFAULT_BOOL",
		"IDCONFIG_TEST_DEFAULT_INT",
		"IDCONFIG_TEST_DEFAULT_INT64",
		"IDCONFIG_TEST_DEFAULT_DURATION",
	} {
		unsetEnv(t, key)
	}

	type config struct {
		Name    string        `env:"IDCONFIG_TEST_DEFAULT_STRING" default:"api"`
		Debug   bool          `env:"IDCONFIG_TEST_DEFAULT_BOOL" default:"true"`
		Port    int           `env:"IDCONFIG_TEST_DEFAULT_INT" default:"8080"`
		Workers int64         `env:"IDCONFIG_TEST_DEFAULT_INT64" default:"4"`
		Timeout time.Duration `env:"IDCONFIG_TEST_DEFAULT_DURATION" default:"30s"`
	}

	var got config
	if err := Load(&got); err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	want := config{Name: "api", Debug: true, Port: 8080, Workers: 4, Timeout: 30 * time.Second}
	if got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadOptionalMissingLeavesExistingValue(t *testing.T) {
	unsetEnv(t, "IDCONFIG_TEST_OPTIONAL_MISSING")

	type config struct {
		Name string `env:"IDCONFIG_TEST_OPTIONAL_MISSING"`
	}

	got := config{Name: "existing"}
	if err := Load(&got); err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if got.Name != "existing" {
		t.Fatalf("Load() changed missing optional value to %q", got.Name)
	}
}

func TestLoadExplicitEmptyString(t *testing.T) {
	t.Setenv("IDCONFIG_TEST_EMPTY_STRING", "")

	type config struct {
		Name string `env:"IDCONFIG_TEST_EMPTY_STRING" default:"default"`
	}

	got := config{Name: "existing"}
	if err := Load(&got); err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if got.Name != "" {
		t.Fatalf("Load() = %q, want an explicitly empty string", got.Name)
	}
}

func TestLoadRequired(t *testing.T) {
	unsetEnv(t, "IDCONFIG_TEST_REQUIRED")

	type config struct {
		Token string `env:"IDCONFIG_TEST_REQUIRED" required:"true"`
	}

	for _, value := range []struct {
		name string
		set  bool
	}{
		{name: "missing"},
		{name: "empty", set: true},
	} {
		t.Run(value.name, func(t *testing.T) {
			if value.set {
				t.Setenv("IDCONFIG_TEST_REQUIRED", "")
			}
			var got config
			err := Load(&got)
			if err == nil || !strings.Contains(err.Error(), "required environment variable is not set") {
				t.Fatalf("Load() error = %v, want required-variable error", err)
			}
		})
	}
}

func TestLoadReportsAllErrorsAndDoesNotMutateTarget(t *testing.T) {
	t.Setenv("IDCONFIG_TEST_PORT", "not-an-int")
	t.Setenv("IDCONFIG_TEST_DEBUG", "not-a-bool")
	t.Setenv("IDCONFIG_TEST_TIMEOUT", "not-a-duration")

	type config struct {
		Port    int           `env:"IDCONFIG_TEST_PORT"`
		Debug   bool          `env:"IDCONFIG_TEST_DEBUG"`
		Timeout time.Duration `env:"IDCONFIG_TEST_TIMEOUT"`
	}

	got := config{Port: 80, Debug: true}
	err := Load(&got)
	if err == nil {
		t.Fatal("Load() returned nil error")
	}
	for _, part := range []string{"IDCONFIG_TEST_PORT", "IDCONFIG_TEST_DEBUG", "IDCONFIG_TEST_TIMEOUT"} {
		if !strings.Contains(err.Error(), part) {
			t.Errorf("Load() error %q does not contain %q", err, part)
		}
	}
	if got != (config{Port: 80, Debug: true}) {
		t.Fatalf("Load() mutated target on error: %+v", got)
	}
}

func TestLoadRejectsInvalidDefault(t *testing.T) {
	unsetEnv(t, "IDCONFIG_TEST_INVALID_DEFAULT")

	type config struct {
		Port int `env:"IDCONFIG_TEST_INVALID_DEFAULT" default:"not-an-int"`
	}

	var got config
	if err := Load(&got); err == nil || !strings.Contains(err.Error(), "must be an integer") {
		t.Fatalf("Load() error = %v, want invalid-default error", err)
	}
}

func TestLoadRequiredFalseAndSkippedField(t *testing.T) {
	unsetEnv(t, "IDCONFIG_TEST_REQUIRED_FALSE")

	type config struct {
		Optional string   `env:"IDCONFIG_TEST_REQUIRED_FALSE" required:"false"`
		Skipped  []string `env:"-"`
		Untagged string
	}

	got := config{Optional: "existing", Skipped: []string{"keep"}, Untagged: "keep"}
	if err := Load(&got); err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if got.Optional != "existing" || len(got.Skipped) != 1 || got.Untagged != "keep" {
		t.Fatalf("Load() changed ignored fields: %+v", got)
	}
}

func TestLoadRejectsInvalidTarget(t *testing.T) {
	tests := []struct {
		name   string
		target any
	}{
		{name: "nil", target: nil},
		{name: "non-pointer", target: struct{}{}},
		{name: "nil pointer", target: (*struct{})(nil)},
		{name: "pointer to non-struct", target: new(int)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := Load(test.target); err == nil {
				t.Fatal("Load() returned nil error")
			}
		})
	}
}

func TestLoadRejectsInvalidTagsAndTypes(t *testing.T) {
	type config struct {
		EmptyEnv    string   `env:""`
		BadRequired string   `env:"IDCONFIG_TEST_BAD_REQUIRED" required:"sometimes"`
		Unsupported []string `env:"IDCONFIG_TEST_UNSUPPORTED" default:"value"`
		private     string   `env:"IDCONFIG_TEST_PRIVATE" default:"value"`
	}

	var got config
	err := Load(&got)
	if err == nil {
		t.Fatal("Load() returned nil error")
	}
	for _, part := range []string{"env tag cannot be empty", "invalid required tag", "unsupported field type", "field must be exported"} {
		if !strings.Contains(err.Error(), part) {
			t.Errorf("Load() error %q does not contain %q", err, part)
		}
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	value, exists := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv(%q) returned an error: %v", key, err)
	}
	t.Cleanup(func() {
		var err error
		if exists {
			err = os.Setenv(key, value)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Errorf("restore environment variable %q: %v", key, err)
		}
	})
}

func TestLoadDoesNotExposeInvalidValue(t *testing.T) {
	const secret = "secret-that-must-not-appear"
	t.Setenv("IDCONFIG_TEST_SECRET_NUMBER", secret)

	type config struct {
		Secret int `env:"IDCONFIG_TEST_SECRET_NUMBER"`
	}

	var got config
	err := Load(&got)
	if err == nil {
		t.Fatal("Load() returned nil error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("Load() exposed the environment value in error %q", err)
	}
}
