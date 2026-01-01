package model

import "testing"

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"abc12345", true},
		{"ABCdef12", true},
		{"short1", false},
		{"abcdefgh", false},
		{"12345678", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := ValidatePasswordStrength(tt.password); got != tt.want {
			t.Fatalf("ValidatePasswordStrength(%q) = %v, want %v", tt.password, got, tt.want)
		}
	}
}
