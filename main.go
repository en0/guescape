package main

import (
	"fmt"
	"os"

	"github.com/en0/guescape/internal/cli"
	"github.com/en0/guescape/internal/hcl"
)

func main() {
	// // TODO: obtain hcl from args
	ws, err := hcl.LoadFile("./docs/example.hcl")
	if err != nil { abort("load error", err, 1) }
	repl := cli.NewRepl(ws, os.Stdin, os.Stdout)
	err = repl.Run()
	if err != nil { abort("game error", err, 1) }
}

func abort(msg string, err error, code int) {
	m := fmt.Sprintf("%s:", msg)
	fmt.Fprintln(os.Stderr, m, err)
	os.Exit(code)
}
