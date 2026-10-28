package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/en0/guescape/internal/game"
)

type uiInterface interface {
    handleAction(game.Action)
    handleQuit()
    handleHelp()
    handleCmdHelp(string)
    handleLook()
    handleInventory()
	handleError(error)
}

type Repl struct {
	running bool
	ws game.WorldState
	out io.Writer
	in *bufio.Scanner
}


func NewRepl(ws game.WorldState, in io.Reader, out io.Writer) Repl {
	bufferedIn := bufio.NewReader(in)
	scanner := bufio.NewScanner(bufferedIn)
	return Repl{ws: ws, in: scanner, out: out}
}

func (r *Repl) handleAction(act game.Action) {
	var ws = r.ws
	for _, ef := range game.PlanEffects(ws, act) {

		switch e := ef.(type) {
		case game.SayEffect:
			fmt.Fprintf(r.out, "%s\n", e.Text)
		case game.EndGameEffect:
			if e.IsWin {
				fmt.Fprintf(r.out, "You Win!\n%s\n", e.Text)
			} else {
				fmt.Fprintf(r.out, "You Lose!\n%s\n", e.Text)
			}
			r.running = false
		default:
			nws, err := game.Apply(ws, e)
			if err != nil {
				r.handleError(err)
				// This is a bug. we could find ourselves in a softlock state. the entire batch
				// should be discarded, not just the one broken effect.
				continue
			}
			ws = nws
		}
	}
	r.ws = ws
}

func (r *Repl) handleQuit() {
	tokens, err := r.prompt("Leave the story here? (y/n) ")
	if err != nil {
		// Treat an error as a clean exit. Likey an EOF
		r.running = false
		return
	}
	if len(tokens) == 1 && strings.HasPrefix("yes", strings.ToLower(tokens[0])) {
		r.running = false
		return
	}

	if len(tokens) == 1 && strings.HasPrefix("no", strings.ToLower(tokens[0])) {
		fmt.Fprint(r.out, "The story continues!\n")
		return
	}

	fmt.Fprint(r.out, "I will take that as a no, the story continues!\n")

}

func (r *Repl) handleHelp() {
	fmt.Fprint(r.out, "available commands:\n")
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	for _, c := range commands {
		fmt.Fprintf(w, " - %s\t%s\n", c.name, c.desc)
	}
	w.Flush()
}

func (r *Repl) handleCmdHelp(s string) {
	cmd, err := findCommand(s)
	if err != nil {
		fmt.Fprintf(r.out, "%q isn't a command. Try \"help\".\n", s)
		return
	}
	fmt.Fprintf(r.out, "%s\n\n", cmd.desc)
	for _, s := range cmd.usage {
		fmt.Fprintf(r.out, "%s\n", s)
	}
}

func (r *Repl) handleLook() {
	v := r.ws.ViewCurrentZone()
	fmt.Fprintf(r.out, "%s\n", v.Desc)
	if len(v.Exits) > 0 {
		fmt.Fprint(r.out, "\nexits:\n")
		for _, e := range v.Exits {
			fmt.Fprintf(r.out, " - %s\n", e.String())
		}
	}
	if len(v.Objects) > 0 {
		fmt.Fprint(r.out, "\nobjects:\n")
		for _, o := range v.Objects {
			fmt.Fprintf(r.out, " - %s\n", o.String())
		}
	}
}

func (r *Repl) handleInventory() {
	inv := r.ws.Inventory()
	if len(inv) > 0 {
		fmt.Fprint(r.out, "objects:\n")
		for _, o := range inv {
			fmt.Fprintf(r.out, " - %s\n", o.String())
		}
	} else {
		fmt.Fprint(r.out, "inventory empty\n")
	}
}

func (r *Repl) handleError(err error) {
	fmt.Fprintf(r.out, "Got an error: %v\n", err)
}

func (r *Repl) prompt(p string) ([]string, error) {
	fmt.Fprint(r.out, p)
	if !r.in.Scan() {
		if err := r.in.Err(); err != nil {
			return nil, fmt.Errorf("reading input: %w", err)
		}
		return nil, io.EOF
	}
	return strings.Fields(r.in.Text()), nil
}

func (r *Repl) Run() error {
	r.running = true
	fmt.Fprintf(r.out, "%s\n", r.ws.Description())
	for r.running {
		v := r.ws.ViewCurrentZone()
		tokens, err := r.prompt(fmt.Sprintf("%s> ", v.Id.String()))
		if errors.Is(err, io.EOF) { return nil }
		if err != nil { return err }
		if len(tokens) == 0 { continue }

		command, err := findCommand(tokens[0])
		if err != nil {
			fmt.Fprintln(r.out, `Those words mean nothing here. Try "help".`)
			continue
		}

		args := tokens[1:]
		if !command.check(args) {
			fmt.Fprintf(
				r.out, "%s doesn't take those words that way. Try \"help %s\".\n",
				command.name, command.name)
			continue
		}

		command.handle(r, args)
	}
	return nil
}
