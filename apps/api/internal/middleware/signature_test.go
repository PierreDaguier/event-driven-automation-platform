package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifyHMACSHA256(t *testing.T) {
	body := []byte(`{"event_type":"lead.created"}`)
	secret := "test-secret"

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if err := VerifyHMACSHA256(sig, body, secret); err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}

	if err := VerifyHMACSHA256("sha256=deadbeef", body, secret); err == nil {
		t.Fatal("expected mismatch error")
	}
}
