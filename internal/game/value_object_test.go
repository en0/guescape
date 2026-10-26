package game

import (
	"errors"
	"testing"
)

func TestNewIdentity(t *testing.T) {
	tests := []struct {
		name string
		input string
		expected Identity
		err error
	} {
		{"valid", "desk-lamp", Identity{value: "desk-lamp"}, nil},
		{"uppercase", "Desk", Identity{}, ErrInvalidIdentity},
		{"space", "desk lamp", Identity{}, ErrInvalidIdentity},
		{"digit", "desk1", Identity{}, ErrInvalidIdentity},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := NewIdentity(tt.input)
			if !errors.Is(err, tt.err) {
				t.Fatalf("NewIdentity(%q) error = %v, want %v", tt.input, err, tt.err)
			}
			if actual != tt.expected || err != tt.err {
				t.Errorf("NewIdentity(%q) = %q want %q", tt.input, actual, tt.expected)
			}
		})
	}
}


