# idtoken — Secure UUID and Token Generators for Go

`idtoken` generates cryptographically secure identifiers and tokens using only
Go's standard library and `crypto/rand`.

## Features

- UUID version 4 with the RFC 9562 version and variant bits
- Lowercase hexadecimal tokens
- Unpadded URL-safe base64 tokens
- Tokens built from a caller-provided printable ASCII alphabet
- No external dependencies

## Installation

```bash
go get github.com/idfactory/idtools@latest
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/idfactory/idtools/idtoken"
)

func main() {
	uuid, err := idtoken.UUID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(uuid)

	hexToken, err := idtoken.Hex(16) // 16 bytes become 32 characters.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hexToken)

	urlToken, err := idtoken.Base64URL(32) // 32 bytes provide 256 bits of entropy.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(urlToken)

	customToken, err := idtoken.String(20, "abcdefghijklmnopqrstuvwxyz0123456789")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(customToken)
}
```

## API contracts

### `UUID() (string, error)`

Returns a UUID version 4 and exposes errors from `crypto/rand` to the caller.

### `Hex(n int) (string, error)`

Generates `n` random bytes and returns `2*n` lowercase hexadecimal characters.
`n` must be positive.

### `Base64URL(n int) (string, error)`

Generates `n` random bytes and encodes them using unpadded URL-safe base64.
`n` must be positive. The result contains no `=`, `+`, or `/` characters.

### `String(length int, alphabet string) (string, error)`

Returns exactly `length` characters selected uniformly from `alphabet`.
`length` must be positive. The alphabet must contain at least two unique
printable ASCII characters from `!` through `~`; spaces, control characters,
Unicode, and duplicates are rejected.

## Error handling

Functions that return an error validate their arguments before reading random
data. Errors from the system random source are wrapped, so callers can inspect
the underlying error with `errors.Is` or `errors.As`.
