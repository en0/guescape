package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/en0/guescape/internal/cli"
	"github.com/en0/guescape/internal/hcl"
)

func main() {
	// // TODO: obtain hcl from args
	ws, err := hcl.LoadFile("./docs/example.hcl")
	if err != nil { abort(err, 1) }
	repl := cli.NewRepl(ws, os.Stdin, os.Stdout)
	err = repl.Run()
	if err != nil { abort(err, 1) }
}

func abort(err error, code int) {
	var loaderError hcl.LoaderError
	if errors.As(err, &loaderError) {
		if errors.Is(err, hcl.ErrParserError) {
			fmt.Fprintf(os.Stderr,
				"Error while parsing %s\n\n%s\n\n",
				loaderError.Context, loaderError.Summary)
		} else if errors.Is(err, hcl.ErrDecodingError) {
			fmt.Fprintf(os.Stderr,
				"Error while decoding %s\n\n%s\n\n",
				loaderError.Context, loaderError.Summary)
		} else {
			fmt.Fprintf(os.Stderr,
				"Unknown error loading %s\n\n%s\n\n",
				loaderError.Context, loaderError.Summary)
		}
		fmt.Fprintln(os.Stderr, loaderError.Detail)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}
