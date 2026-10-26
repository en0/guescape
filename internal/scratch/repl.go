package scratch

import (
	"bufio"
	"fmt"
	"io"

	"github.com/en0/guescape/internal/game"
)

func Play(ws game.WorldState, input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	for {
		// collect user input
		if !scanner.Scan() {
			break
		}

		input := tokenizeInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		args := input[1:]
		cmd, err := normCmd(input[0])
		if err != nil {
			fmt.Fprintln(output, err); continue
		}

		fmt.Fprintln(output, "CMD:", cmd.cmd)
		fmt.Fprintln(output, "ARGS:", args)
		fmt.Fprintln(output, "CLI:", cmd.cli)
		fmt.Fprintln(output, "CHECK:", cmd.check(args))
		fmt.Fprintln(output, "DESC:", cmd.desc)
		fmt.Fprintln(output, "USAGE:", cmd.usage)

		if !cmd.check(args) {
			fmt.Fprintf(output, "%s what?\n", cmd.cmd)
		}

		// create game.Action
		// collect effects
		// apply effects
		// render feedback
	}
	return nil
}
