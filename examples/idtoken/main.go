package main

import (
	"fmt"

	"github.com/idfactory/idtools/idlog"
	"github.com/idfactory/idtools/idtoken"
)

func main() {
	uuid, err := idtoken.UUID()
	if err != nil {
		idlog.Fail(fmt.Errorf("generate UUID: %w", err))
	}
	fmt.Printf("UUID: %s\n", uuid)

	hexToken, err := idtoken.Hex(16)
	if err != nil {
		idlog.Fail(fmt.Errorf("generate hex token: %w", err))
	}
	fmt.Printf("Hex token: %s\n", hexToken)

	urlToken, err := idtoken.Base64URL(32)
	if err != nil {
		idlog.Fail(fmt.Errorf("generate URL-safe token: %w", err))
	}
	fmt.Printf("URL-safe token: %s\n", urlToken)

	customToken, err := idtoken.String(20, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if err != nil {
		idlog.Fail(fmt.Errorf("generate custom token: %w", err))
	}
	fmt.Printf("Custom token: %s\n", customToken)
}
