package cli

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/en0/guescape/internal/game"
)

func testIdentity(t *testing.T, id string) game.Identity {
	t.Helper()
	ident, err := game.NewIdentity(id)
	if err != nil { t.Fatalf("game.NewIdent(%q) err=%v", id, err) }
	return ident
}

func testRepl(t *testing.T) (Repl, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	state, err := game.NewWorldState(
		testIdentity(t, "zone-one"),
		"unittest-world",
		[]game.Identity{},
		[]game.Zone{
			game.NewZone(testIdentity(t, "zone-one"), "zone-one", nil),
		},
		[]game.Object{},
		nil, nil,
	)
	if err != nil { t.Fatalf("game.NewWorldState(...) err=%v", err) }
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	return NewRepl(state, in, out), in, out
}

func TestReplPrompt(t *testing.T) {

	tests := []struct{
		in string
		out []string
		err error
	}{
		{"word", []string{"word"}, nil},
		{"wo-rd", []string{"wo-rd"}, nil},
		{"wo\trd", []string{"wo", "rd"}, nil},
		{"wo   rd", []string{"wo", "rd"}, nil},
		{"  word  ", []string{"word"}, nil},
		{"\t word \t", []string{"word"}, nil},
		{"CA PS", []string{"CA", "PS"}, nil},
		{"12 34", []string{"12", "34"}, nil},
		{"§ §", []string{"§", "§"}, nil},

		{"\n", []string{}, nil},
		{"", []string{}, io.EOF},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			repl, inBuffer, _ := testRepl(t)
			inBuffer.WriteString(tt.in)
			result, err := repl.prompt("")
			if !errors.Is(err, tt.err) {
				t.Fatalf("repl.prompt(%q) err=%v, want err=%v", tt.in, err, tt.err)
			}
			if !slices.Equal(result, tt.out) {
				t.Fatalf("repl.prompt(%q) %v, want %v", tt.in, result, tt.out)
			}
		})
	}
}

// r, _, out := testRepl(t, "look north\n")
//
// tokens, err := r.prompt("> ")
// if err != nil {
//     t.Fatalf("prompt() err=%v", err)
// }
//
// want := []string{"look", "north"}
// if !reflect.DeepEqual(tokens, want) {
//     t.Errorf("tokens = %v, want %v", tokens, want)
// }
//
// if got := out.String(); got != "> " {
//     t.Errorf("prompt text = %q, want %q", got, "> ")
// }
