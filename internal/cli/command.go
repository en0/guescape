package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/en0/guescape/internal/game"
)

type command struct {
	name string
	desc string
	usage []string
	alias []string
	check func([]string) bool
	handle func(uiInterface, []string)
}

var commands = []command{
    {
		name: "help",
		desc: "Show help for all commands, or one command",
		usage: []string{
			"help [command]",
			"  command  the command to show help for",
		},
		alias: []string{"?"},
		check: func(s []string) bool {a := len(s); return a == 0 || a == 1},
		handle: func(u uiInterface, s []string) {
			if len(s) == 0 { u.handleHelp() }
			if len(s) == 1 { u.handleCmdHelp(s[0]) }
		},
	},
    {
		name: "quit",
		alias: []string{"exit"},
		desc: "Exit the game",
		usage: []string{
			"quit",
		},
		check: func(s []string) bool {return len(s) == 0},
		handle: func(u uiInterface, s []string) { u.handleQuit() },
	},
    {
		name: "inventory",
		desc: "List the objects you're carrying",
		usage: []string{
			"inventory",
		},
		check: func(s []string) bool {return len(s) == 0},
		handle: func(u uiInterface, s []string) { u.handleInventory() },
	},
    {
		name: "look",
		alias: []string{"show"},
		desc: "Look around",
		usage: []string{
			"look",
		},
		check: func(s []string) bool {return len(s) == 0},
		handle: func(u uiInterface, s []string) { u.handleLook() },
	},

    // Passed to the game.
    {
		name: "move",
		desc: "Move in a direction",
		usage: []string{
			"move <direction>",
			"  direction  the direction to move",
		},
		check: func(s []string) bool {return len(s) == 1},
		handle: func(u uiInterface, s []string) {
			id, err := game.NewIdentity(strings.ToLower(s[0]))
			if err != nil { u.handleError(err); return }
			act := game.NewNavAction(game.CmdMove, id)
			u.handleAction(act)
		},
	},
    {
		name: "examine",
		alias: []string{"inspect"},
		desc: "Examine an object",
		usage: []string{
			"examine <object>",
			"  object  the object to examine",
		},
		check: func(s []string) bool {return len(s) == 1},
		handle: func(u uiInterface, s []string) {
			id, err := game.NewIdentity(strings.ToLower(s[0]))
			if err != nil { u.handleError(err); return }
			act := game.NewObjectAction(game.VerbExamine, id)
			u.handleAction(act)
		},
	},
    {
		name: "use",
		desc: "Use an object, optionally on another",
		usage: []string{
			"use <object> [on <target>]",
			"  object  the object to use",
			"  target  what to use it on, if any",
		},
		check: func(s []string) bool {a := len(s); return a == 1 || a == 3},
		handle: func(u uiInterface, s []string) {
			if len(s) == 1 {
				id, err := game.NewIdentity(strings.ToLower(s[0]))
				if err != nil { u.handleError(err); return }
				act := game.NewObjectAction(game.VerbUse, id)
				u.handleAction(act)
			}
			if len(s) == 3 {
				iid, err := game.NewIdentity(strings.ToLower(s[0]))
				if err != nil { u.handleError(err); return }
				oid, err := game.NewIdentity(strings.ToLower(s[2]))
				if err != nil { u.handleError(err); return }
				act := game.NewObjectActionWithItem(game.VerbUse, oid, iid)
				u.handleAction(act)
			}
		},
	},
    {
		name: "take",
		desc: "Pick up an object",
		usage: []string{
			"take <object>",
			"  object  the object to take",
		},
		check: func(s []string) bool {return len(s) == 1},
		handle: func(u uiInterface, s []string) {
			id, err := game.NewIdentity(strings.ToLower(s[0]))
			if err != nil { u.handleError(err); return }
			act := game.NewObjectAction(game.VerbTake, id)
			u.handleAction(act)
		},
	},
    {
		name: "drop",
		desc: "Put down an object",
		usage: []string{
			"drop <object>",
			"  object  the object to drop",
		},
		check: func(s []string) bool {return len(s) == 1},
		handle: func(u uiInterface, s []string) {
			id, err := game.NewIdentity(strings.ToLower(s[0]))
			if err != nil { u.handleError(err); return }
			act := game.NewObjectAction(game.VerbDrop, id)
			u.handleAction(act)
		},
	},
}

var errUnknownCommand = errors.New("unknown command")

func findCommand(cmd string) (command, error) {
	cmd = strings.ToLower(cmd)
	for _, c := range commands {
		if strings.HasPrefix(c.name, cmd) {
			return c, nil
		} else if slices.Contains(c.alias, cmd) {
			return c, nil
		}
	}
	return command{}, fmt.Errorf("%w: %q", errUnknownCommand, cmd)
}
