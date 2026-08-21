package idtoken

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"
)

var errRandomSource = errors.New("random source failed")

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errRandomSource
}

func TestUUID(t *testing.T) {
	value, err := UUID()
	if err != nil {
		t.Fatalf("UUID() returned an error: %v", err)
	}
	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !pattern.MatchString(value) {
		t.Fatalf("UUID() returned invalid UUID v4 %q", value)
	}
}

func TestUUIDWithReader(t *testing.T) {
	input := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	value, err := uuidWithReader(strings.NewReader(string(input)))
	if err != nil {
		t.Fatalf("uuidWithReader() returned an error: %v", err)
	}
	const want = "00010203-0405-4607-8809-0a0b0c0d0e0f"
	if value != want {
		t.Fatalf("uuidWithReader() = %q, want %q", value, want)
	}
}

func TestUUIDWithReaderError(t *testing.T) {
	_, err := uuidWithReader(failingReader{})
	if !errors.Is(err, errRandomSource) {
		t.Fatalf("uuidWithReader() error = %v, want wrapped source error", err)
	}
}

func TestHexWithReader(t *testing.T) {
	value, err := hexWithReader(strings.NewReader("\x00\xab\xff"), 3)
	if err != nil {
		t.Fatalf("hexWithReader() returned an error: %v", err)
	}
	if value != "00abff" {
		t.Fatalf("hexWithReader() = %q, want %q", value, "00abff")
	}
}

func TestBase64URLWithReader(t *testing.T) {
	value, err := base64URLWithReader(strings.NewReader("\xfb\xff"), 2)
	if err != nil {
		t.Fatalf("base64URLWithReader() returned an error: %v", err)
	}
	if value != "-_8" {
		t.Fatalf("base64URLWithReader() = %q, want %q", value, "-_8")
	}
	if strings.ContainsAny(value, "=+/") {
		t.Fatalf("base64URLWithReader() returned non-URL-safe value %q", value)
	}
}

func TestByteGeneratorsErrors(t *testing.T) {
	tests := []struct {
		name     string
		generate func(io.Reader, int) (string, error)
	}{
		{name: "hex", generate: hexWithReader},
		{name: "base64 URL", generate: base64URLWithReader},
	}

	for _, test := range tests {
		t.Run(test.name+" invalid size", func(t *testing.T) {
			for _, n := range []int{-1, 0} {
				if _, err := test.generate(strings.NewReader(""), n); err == nil {
					t.Fatalf("generator with size %d returned nil error", n)
				}
			}
		})

		t.Run(test.name+" source failure", func(t *testing.T) {
			_, err := test.generate(failingReader{}, 1)
			if !errors.Is(err, errRandomSource) {
				t.Fatalf("generator error = %v, want wrapped source error", err)
			}
		})
	}
}

func TestStringWithReader(t *testing.T) {
	// 255 is rejected for a three-character alphabet. The remaining values
	// map deterministically to A, B, and C.
	value, err := stringWithReader(strings.NewReader("\xff\x00\x01\x02"), 3, "ABC")
	if err != nil {
		t.Fatalf("stringWithReader() returned an error: %v", err)
	}
	if value != "ABC" {
		t.Fatalf("stringWithReader() = %q, want %q", value, "ABC")
	}
}

func TestString(t *testing.T) {
	const alphabet = "abcXYZ012"
	value, err := String(64, alphabet)
	if err != nil {
		t.Fatalf("String() returned an error: %v", err)
	}
	if len(value) != 64 {
		t.Fatalf("len(String()) = %d, want 64", len(value))
	}
	for _, character := range value {
		if !strings.ContainsRune(alphabet, character) {
			t.Fatalf("String() returned character %q outside alphabet", character)
		}
	}
}

func TestStringValidation(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		alphabet string
	}{
		{name: "negative length", length: -1, alphabet: "ab"},
		{name: "zero length", length: 0, alphabet: "ab"},
		{name: "empty alphabet", length: 1, alphabet: ""},
		{name: "single character", length: 1, alphabet: "a"},
		{name: "duplicate character", length: 1, alphabet: "aba"},
		{name: "space", length: 1, alphabet: "a b"},
		{name: "control character", length: 1, alphabet: "a\nb"},
		{name: "unicode", length: 1, alphabet: "aя"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := stringWithReader(strings.NewReader(""), test.length, test.alphabet); err == nil {
				t.Fatal("stringWithReader() returned nil error")
			}
		})
	}
}

func TestStringSourceError(t *testing.T) {
	_, err := stringWithReader(failingReader{}, 1, "ab")
	if !errors.Is(err, errRandomSource) {
		t.Fatalf("stringWithReader() error = %v, want wrapped source error", err)
	}
}
