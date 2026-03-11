package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

func VerifyHMACSHA256(signatureHeader string, body []byte, secret string) error {
	if secret == "" {
		return errors.New("missing webhook secret")
	}
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return errors.New("invalid signature format")
	}
	expectedMAC := hmac.New(sha256.New, []byte(secret))
	expectedMAC.Write(body)
	expected := expectedMAC.Sum(nil)

	got, err := hex.DecodeString(parts[1])
	if err != nil {
		return errors.New("invalid signature encoding")
	}
	if !hmac.Equal(expected, got) {
		return errors.New("signature mismatch")
	}
	return nil
}
