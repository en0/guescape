package scratch

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func testCommand(t *testing.T, cmd string) command {
	t.Helper()
	i := slices.IndexFunc(commands, func(c command) bool {
		return c.cmd == cmd
	})
	if i < 0 { t.Fatalf("testCommand: bad command %q", cmd) }
	return commands[i]
}

func TestNormCmd(t *testing.T) {

	tests := []struct {
		input string
		expect command
		err error
	} {
		// Short codes
		{input: "h", expect: testCommand(t, "help")},
		{input: "q", expect: testCommand(t, "quit")},
		{input: "l", expect: testCommand(t, "look")},
		{input: "m", expect: testCommand(t, "move")},
		{input: "i", expect: testCommand(t, "inventory")},
		{input: "e", expect: testCommand(t, "examine")},
		{input: "t", expect: testCommand(t, "take")},
		{input: "u", expect: testCommand(t, "use")},
		{input: "d", expect: testCommand(t, "drop")},

		// Aliases
		{input: "?", 		expect: testCommand(t, "help")},
		{input: "show", 	expect: testCommand(t, "look")},
		{input: "inspect", 	expect: testCommand(t, "examine")},
		{input: "exit", 	expect: testCommand(t, "quit")},

		// Errors
		{input: "no-exist", err: errUnknownCommand},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := normCmd(tt.input)
			if !errors.Is(err, tt.err) {
				t.Fatalf("normCmd(%q) err=%v", tt.input, err)
			}
			if actual.cmd != tt.expect.cmd {
				t.Errorf(
					"normCmd(%q) = (%v), want (%v)",
					tt.input, actual, tt.expect,
				)
			}
		})
	}

}

func TestTokenizeInput(t *testing.T) {
	tests := []struct {
		input string
		expected []string
	} {
		{"hello world", []string{"hello", "world"}},
		{"HELLO WORLD", []string{"hello", "world"}},
		{"This is a Test", []string{"this", "is", "a", "test"}},
		{"command", []string{"command"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := tokenizeInput(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf(
					"tokenizeInput(%q) = %q, want %q",
					tt.input, actual, tt.expected,
				)
			}
		})
	}
}

func TestCommandUniquePrefix(t *testing.T) {
	for i, a := range commands {
		for j, b := range commands {
			if i == j { continue }
			name := fmt.Sprintf("%v-%v", a.cmd, b.cmd)
			t.Run(name, func(t *testing.T) {
				if strings.HasPrefix(a.cmd, b.cmd[:1]) {
					t.Fatalf("%q and %q share a prefix", a.cmd, b.cmd)
				}
			})
		}
	}
}
