package game

import "testing"

 func TestMatchIdentity(t *testing.T) {
	 tests := []struct {
		 name string
		 partial string
		 candidates []Identity
		 expectIdentity Identity
		 expectFound bool
	 }{
		 {
			name: "no-match-returns-unsuccessfull",
			partial: "xyz",
			candidates: testIdentities(t, "abcd", "efgh", "ijkl"),
			expectFound: false,
		 },
		 {
			name: "matches-on-prefix-returns-identity",
			partial: "e",
			candidates: testIdentities(t, "abcd", "efgh", "ijkl"),
			expectIdentity: testIdentity(t, "efgh"),
			expectFound: true,
		 },
		 {
			name: "matches-on-suffix-returns-identity",
			partial: "h",
			candidates: testIdentities(t, "abcd", "efgh", "ijkl"),
			expectIdentity: testIdentity(t, "efgh"),
			expectFound: true,
		 },
		 {
			name: "matches-on-middle-returns-identity",
			partial: "g",
			candidates: testIdentities(t, "abcd", "efgh", "ijkl"),
			expectIdentity: testIdentity(t, "efgh"),
			expectFound: true,
		 },
		 {
			name: "matches-multiple-returns-unsuccessfull",
			partial: "ab",
			candidates: testIdentities(t, "abc", "abd"),
			expectFound: false,
		 },
		 {
			name: "matches-case-insensitive",
			partial: "EFGH",
			candidates: testIdentities(t, "abcd", "efgh", "ijkl"),
			expectIdentity: testIdentity(t, "efgh"),
			expectFound: true,
		 },
		 {
			name: "empty-string-matches-nothing",
			partial: "",
			candidates: testIdentities(t, "abcd"),
			expectFound: false,
		 },
	 }

	 for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			 identity, found := matchIdentity(tt.candidates, tt.partial)
			 if found != tt.expectFound || identity != tt.expectIdentity {
				 t.Fatalf(
					 "MatchIdentity(%q) (%v, %v), want (%v, %v)",
					 tt.partial, identity, found,
					 tt.expectIdentity, tt.expectFound)
			 }
		 })
	 }
 }
