package service

import "testing"

func TestNormalizeCreateAPIKeyInput(t *testing.T) {
	rate, platforms := normalizeCreateAPIKeyInput(&CreateAPIKeyRequest{
		RateLimit:        0,
		AllowedPlatforms: "",
	})
	if rate != 0 {
		t.Fatalf("expected rate limit 0, got %d", rate)
	}
	if platforms != "all" {
		t.Fatalf("expected platforms all, got %q", platforms)
	}

	rate, platforms = normalizeCreateAPIKeyInput(&CreateAPIKeyRequest{
		RateLimit:        -5,
		AllowedPlatforms: "openai",
	})
	if rate != 0 {
		t.Fatalf("expected rate limit 0 for negative, got %d", rate)
	}
	if platforms != "openai" {
		t.Fatalf("expected platforms openai, got %q", platforms)
	}
}
