package service

import (
	"testing"
	"time"
)

func TestMemoryStoreVerifyAndExpire(t *testing.T) {
	store := NewMemoryStore(20 * time.Millisecond)

	if err := store.Set("id1", "abcd"); err != nil {
		t.Fatalf("set captcha failed: %v", err)
	}

	if got := store.Get("id1", false); got != "abcd" {
		t.Fatalf("expected captcha value, got %q", got)
	}

	if ok := store.Verify("id1", "abcd", true); !ok {
		t.Fatal("expected verify success")
	}

	if got := store.Get("id1", false); got != "" {
		t.Fatalf("expected captcha cleared, got %q", got)
	}

	if err := store.Set("id2", "efgh"); err != nil {
		t.Fatalf("set captcha failed: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if ok := store.Verify("id2", "efgh", false); ok {
		t.Fatal("expected captcha to be expired")
	}
}
