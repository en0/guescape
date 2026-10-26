package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/en0/guescape/internal/game"
	"github.com/en0/guescape/internal/hcl"
)

func main2() {
	ws, err := hcl.LoadFile("./docs/example.hcl")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load error:", err)
		os.Exit(1)
	}

	fmt.Println("guescape — type 'quit' to exit")
	fmt.Println(strings.Repeat("-", 40))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // EOF (Ctrl-D) or read error
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "quit" || line == "q" || line == "exit" {
			break
		}

		action, ok := parseAction(line)
		if !ok {
			fmt.Println("I don't understand that.")
			continue
		}

		// domain: plan the effects, then route them.
		for _, ef := range game.PlanEffects(ws, action) {
			switch e := ef.(type) {
			case game.SayEffect:
				fmt.Println(e.Text)
			case game.EndGameEffect:
				fmt.Println(e.Text)
				if e.IsWin {
					fmt.Println("** You win. **")
				} else {
					fmt.Println("** Game over. **")
				}
				return
			default:
				next, err := game.Apply(ws, ef)
				if err != nil {
					// Contract: discard the result, keep the old state.
					continue
				}
				ws = next
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
	fmt.Println("bye")
}

// parseAction turns a raw line into a game.Action. Nasty first pass:
//   move <exit>
//   examine|take|use|drop <object>
//   use <object> on <item>
func parseAction(line string) (game.Action, bool) {
	fields := strings.Fields(strings.ToLower(line))
	if len(fields) == 0 {
		return game.Action{}, false
	}

	verb := resolveAlias(fields[0])
	args := fields[1:]

	switch verb {
	case "move":
		if len(args) != 1 {
			return game.Action{}, false
		}
		exit, err := game.NewIdentity(args[0])
		if err != nil {
			return game.Action{}, false
		}
		return game.NewNavAction(game.CmdMove, exit), true

	case "examine", "take", "use", "drop":
		if len(args) == 0 {
			return game.Action{}, false
		}
		objVerb, ok := objectVerb(verb)
		if !ok {
			return game.Action{}, false
		}

		// "use <item> on <target>": item is used with the target.
		if verb == "use" && len(args) == 3 && args[1] == "on" {
			item, err := game.NewIdentity(args[0])
			if err != nil {
				return game.Action{}, false
			}
			target, err := game.NewIdentity(args[2])
			if err != nil {
				return game.Action{}, false
			}
			return game.NewObjectActionWithItem(objVerb, target, item), true
		}

		// "<verb> <target>"
		target, err := game.NewIdentity(args[0])
		if err != nil {
			return game.Action{}, false
		}
		return game.NewObjectAction(objVerb, target), true
	}

	return game.Action{}, false
}

func objectVerb(s string) (game.ObjectVerb, bool) {
	switch s {
	case "examine":
		return game.VerbExamine, true
	case "take":
		return game.VerbTake, true
	case "use":
		return game.VerbUse, true
	case "drop":
		return game.VerbDrop, true
	}
	return 0, false
}

func resolveAlias(raw string) string {
	aliases := map[string]string{
		"m": "move",
		"e": "examine",
		"inspect": "examine",
		"t": "take",
		"u": "use",
		"d": "drop",
	}
	if cmd, ok := aliases[raw]; ok {
		return cmd
	}
	return raw
}
