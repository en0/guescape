package game

import (
	"fmt"
	"strings"
)

// A cross-entity pointer that doesn't resolve
type ReferenceError struct {
	kind error
	// The entity holding a bad reference
	Referrer Identity
	// The missing thing being referred to
	Referent Identity
	// The container the referent should belong to
	Scope Identity
}

func (e ReferenceError) Unwrap() error { return e.kind }
func (e ReferenceError) Error() string {
	var msg []string
	if e.Referrer != (Identity{}) {
		msg = append(msg, fmt.Sprintf("referrer=%q", e.Referrer))
	}
	if e.Referent != (Identity{}) {
		msg = append(msg, fmt.Sprintf("referent=%q", e.Referent))
	}
	if e.Scope != (Identity{}) {
		msg = append(msg, fmt.Sprintf("scope=%q", e.Scope))
	}
	if len(msg) > 0 {
		return fmt.Sprintf("%s: %s", e.kind, strings.Join(msg, ", "))
	}
	return e.kind.Error()
}
