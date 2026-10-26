package hcl

import (
	"github.com/en0/guescape/internal/game"
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

// TODO: Make these carry the right information
// Right now we are dropping a ton of details that we want

func parsingErrorMapper(diag hcl.Diagnostic) error {
	return game.NewParsingError(diag.Summary)
}

func decodingErrorMapper(diag hcl.Diagnostic) error {
	return game.NewDecodingError(diag.Summary)
}
