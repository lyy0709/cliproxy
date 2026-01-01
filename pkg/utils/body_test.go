package utils

import (
	"bytes"
	"testing"
)

func TestReadAllWithLimit(t *testing.T) {
	data := bytes.Repeat([]byte("a"), 10)
	_, err := ReadAllWithLimit(bytes.NewReader(data), 5)
	if err != ErrBodyTooLarge {
		t.Fatalf("expected ErrBodyTooLarge, got %v", err)
	}

	limited, err := ReadAllWithLimit(bytes.NewReader(data), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(limited) != 10 {
		t.Fatalf("expected length 10, got %d", len(limited))
	}
}
