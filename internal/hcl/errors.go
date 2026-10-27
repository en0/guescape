package hcl

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type diagMapper func(diag hcl.Diagnostic) error

func appendDiagsAs(errs []error, diags hcl.Diagnostics, remap diagMapper) []error {
	if !diags.HasErrors() {
		return errs
	}
	for _, d := range diags {
		if d.Severity == hcl.DiagError {
			errs = append(errs, remap(*d))
		}
	}
	return errs
}

var ErrParserError = errors.New("parser error")
var ErrDecodingError = errors.New("decoding error")

type LoaderError struct{
	kind error
	Summary string
	Detail string
	Context string
}

func (e LoaderError) Unwrap() error { return e.kind }
func (e LoaderError) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.Summary)
}

func parsingErrorMapper(diag hcl.Diagnostic) error {
	return LoaderError{
		kind: ErrParserError,
		Summary: diag.Summary,
		Detail: diag.Detail,
		Context: diag.Subject.String(),
	}
}

func decodingErrorMapper(diag hcl.Diagnostic) error {
	return LoaderError{
		kind: ErrDecodingError,
		Summary: diag.Summary,
		Detail: diag.Detail,
		Context: diag.Subject.String(),
	}
}
