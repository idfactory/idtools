// Package idtoken provides cryptographically secure UUID and token generators.
package idtoken

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

const uuidSize = 16

// UUID returns a randomly generated UUID version 4.
// It returns an error if the system's cryptographically secure random source
// cannot provide enough data.
func UUID() (string, error) {
	return uuidWithReader(rand.Reader)
}

// Hex returns n cryptographically secure random bytes encoded as 2*n
// lowercase hexadecimal characters. It returns an error when n is not
// positive or the random source fails.
func Hex(n int) (string, error) {
	return hexWithReader(rand.Reader, n)
}

// Base64URL returns n cryptographically secure random bytes encoded with
// unpadded URL-safe base64. It returns an error when n is not positive or the
// random source fails.
func Base64URL(n int) (string, error) {
	return base64URLWithReader(rand.Reader, n)
}

// String returns a cryptographically secure random string of the requested
// length using alphabet. The alphabet must contain at least two unique
// printable ASCII characters in the range '!' through '~'.
func String(length int, alphabet string) (string, error) {
	return stringWithReader(rand.Reader, length, alphabet)
}

func uuidWithReader(reader io.Reader) (string, error) {
	b, err := readRandom(reader, uuidSize)
	if err != nil {
		return "", err
	}

	// Set the UUID version to 4 and the variant to RFC 9562.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func hexWithReader(reader io.Reader, n int) (string, error) {
	if err := validateByteCount(n); err != nil {
		return "", err
	}

	b, err := readRandom(reader, n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func base64URLWithReader(reader io.Reader, n int) (string, error) {
	if err := validateByteCount(n); err != nil {
		return "", err
	}

	b, err := readRandom(reader, n)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func stringWithReader(reader io.Reader, length int, alphabet string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("idtoken: length must be positive: %d", length)
	}
	if err := validateAlphabet(alphabet); err != nil {
		return "", err
	}

	result := make([]byte, length)
	bufferSize := min(length, 128)
	buffer := make([]byte, bufferSize)
	limit := 256 - (256 % len(alphabet))

	for position := 0; position < length; {
		remaining := length - position
		batch := min(remaining, len(buffer))
		if _, err := io.ReadFull(reader, buffer[:batch]); err != nil {
			return "", fmt.Errorf("idtoken: read random bytes: %w", err)
		}

		for _, value := range buffer[:batch] {
			if int(value) >= limit {
				continue
			}
			result[position] = alphabet[int(value)%len(alphabet)]
			position++
			if position == length {
				break
			}
		}
	}

	return string(result), nil
}

func validateByteCount(n int) error {
	if n <= 0 {
		return fmt.Errorf("idtoken: byte count must be positive: %d", n)
	}
	return nil
}

func validateAlphabet(alphabet string) error {
	if len(alphabet) < 2 {
		return fmt.Errorf("idtoken: alphabet must contain at least two characters")
	}

	var seen [128]bool
	for i := 0; i < len(alphabet); i++ {
		character := alphabet[i]
		if character < '!' || character > '~' {
			return fmt.Errorf("idtoken: alphabet must contain only printable ASCII characters without spaces")
		}
		if seen[character] {
			return fmt.Errorf("idtoken: alphabet contains duplicate character %q", character)
		}
		seen[character] = true
	}
	return nil
}

func readRandom(reader io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(reader, b); err != nil {
		return nil, fmt.Errorf("idtoken: read random bytes: %w", err)
	}
	return b, nil
}
