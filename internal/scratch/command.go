package scratch

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func tokenizeInput(input string) ([]string) {
	input = strings.TrimSpace(input)
	if input == "" { return []string{} }
	result := strings.Split(input, " ")
	for i, token := range result {
		v := strings.TrimSpace(token)
		result[i] = strings.ToLower(v)
	}
	return result
}

var errUnknownCommand = errors.New("unknown command")

type command struct {
    cmd  string
	alias []string
    usage string
    desc  string
	check func ([]string) bool
    cli   bool
}

var commands = []command{
    // Handled by the CLI.
    {
		cmd: "help",
		alias: []string{"?"},
		usage: "help [cmd]",
		desc: "Display this message, or help for one command",
		check: func(s []string) bool {a := len(s); return a == 0 || a == 1},
		cli: true,
	},
    {
		cmd: "quit",
		alias: []string{"exit"},
		usage: "quit",
		desc: "Exit the game",
		check: func(s []string) bool {return len(s) == 0},
		cli: true,
	},
    {
		cmd: "inventory",
		usage: "inventory",
		desc: "Look at your inventory",
		check: func(s []string) bool {return len(s) == 0},
		cli: true,
	},
    {
		cmd: "look",
		alias: []string{"show"},
		usage: "look",
		desc: "Look around",
		check: func(s []string) bool {return len(s) == 0},
		cli: true,
	},

    // Passed to the game.
    {
		cmd: "move",
		usage: "move <direction>",
		desc: "Move in a direction",
		check: func(s []string) bool {return len(s) == 1},
		cli: false,
	},
    {
		cmd: "examine",
		alias: []string{"inspect"},
		usage: "examine <object>",
		desc: "Examine an object",
		check: func(s []string) bool {return len(s) == 1},
		cli: false,
	},
    {
		cmd: "take",
		usage: "take <object>",
		desc: "Move an object to your inventory",
		check: func(s []string) bool {return len(s) == 1},
		cli: false,
	},
    {
		cmd: "use",
		usage: "use <object> [on <target>]",
		desc: "Use an object, optionally on another",
		check: func(s []string) bool {a := len(s); return a == 1 || a == 3},
		cli: false,
	},
    {
		cmd: "drop",
		usage: "drop <object>",
		desc: "Remove an object from your inventory",
		check: func(s []string) bool {return len(s) == 1},
		cli: false,
	},
}

func normCmd(cmd string) (command, error) {
	for _, c := range commands {
		if strings.HasPrefix(c.cmd, cmd) {
			return c, nil
		} else if slices.Contains(c.alias, cmd) {
			return c, nil
		}
	}
	return command{}, fmt.Errorf("%w: %q", errUnknownCommand, cmd)
}
