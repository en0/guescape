package game

import (
	"fmt"
	"strings"
)

// SourceError: Returned by Loader.
// indicates the input file could not be read.
type SourceError struct {
	err error
}
func (e SourceError) Error() string {
	return e.err.Error()
}
func (e SourceError) Unwrap() error {
	return e.err
}
func NewSourceError (err error) SourceError {
	return SourceError{err}
}

// ParsingError: Returned by the loader
// indicates an invalid syntax in the input file.
type ParsingError struct {
	msg string
}
func (e ParsingError) Error() string {
	return e.msg
}
func NewParsingError (msg string) ParsingError {
	return ParsingError{
		msg: msg,
	}
}

// DecodingError: Returned by the loader
// indicates the input file could not be decoded into the expected structure.
type DecodingError struct {
	msg string
}
func (e DecodingError) Error() string {
	return e.msg
}
func NewDecodingError(msg string) DecodingError {
	return DecodingError{
		msg: msg,
	}
}

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
