package utils

import (
	"strings"
	"testing"

	"cli-proxy/internal/config"
)

func TestEncryptDecryptString(t *testing.T) {
	origCfg := config.Cfg
	config.Cfg = &config.Config{
		Security: config.SecurityConfig{DataKey: "unit-test-key"},
		JWT:      config.JWTConfig{Secret: "fallback-key"},
	}
	defer func() {
		config.Cfg = origCfg
	}()

	plain := "secret-value-123"
	encrypted, err := EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString error: %v", err)
	}
	if encrypted == plain {
		t.Fatalf("expected encrypted value, got plaintext")
	}
	if !strings.HasPrefix(encrypted, "enc:") {
		t.Fatalf("expected encrypted prefix, got %q", encrypted)
	}

	decrypted, err := DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString error: %v", err)
	}
	if decrypted != plain {
		t.Fatalf("expected %q, got %q", plain, decrypted)
	}

	noop, err := DecryptString(plain)
	if err != nil {
		t.Fatalf("DecryptString(plain) error: %v", err)
	}
	if noop != plain {
		t.Fatalf("expected plaintext passthrough, got %q", noop)
	}
}
