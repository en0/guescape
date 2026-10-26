package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/en0/guescape/internal/game"
)

type testCallResult[T any] struct {
	called bool
	calledWith T
}

type testUiSpy struct {
    handleActionCalled testCallResult[game.Action]
    handleQuitCalled testCallResult[int]
    handleHelpCalled testCallResult[int]
    handleCmdHelpCalled testCallResult[string]
    handleLookCalled testCallResult[int]
    handleInventoryCalled testCallResult[int]
	handleErrorCalled testCallResult[error]
}

func (t *testUiSpy) handleAction(a game.Action) {
	t.handleActionCalled = testCallResult[game.Action]{true, a}
}

func (t *testUiSpy) handleQuit() {
	t.handleQuitCalled = testCallResult[int]{true, 0}
}

func (t *testUiSpy) handleHelp() {
	t.handleHelpCalled = testCallResult[int]{true, 0}
}

func (t *testUiSpy) handleCmdHelp(s string) {
	t.handleCmdHelpCalled = testCallResult[string]{true, s}
}

func (t *testUiSpy) handleLook() {
	t.handleLookCalled = testCallResult[int]{true, 0}
}

func (t *testUiSpy) handleInventory() {
	t.handleInventoryCalled = testCallResult[int]{true, 0}
}

func (t *testUiSpy) handleError(e error) {
	t.handleErrorCalled = testCallResult[error]{true, e}
}

func testCommand(t *testing.T, name string) command {
	t.Helper()
	i := slices.IndexFunc(commands, func(c command) bool {
		return c.name == name
	})
	if i < 0 { t.Fatalf("testCommand: bad command %q", name) }
	return commands[i]
}

func TestCommandUniquePrefix(t *testing.T) {
	for i, a := range commands {
		for j, b := range commands {
			if i == j { continue }
			name := fmt.Sprintf("%v-%v", a.name, b.name)
			t.Run(name, func(t *testing.T) {
				if strings.HasPrefix(a.name, b.name[:1]) {
					t.Fatalf("commands[%v] %q and %q share a prefix", i, a.name, b.name)
				}
			})
		}
	}
}

func TestFindCommand(t *testing.T) {

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

		// misc
		{input: "help", expect: testCommand(t, "help")},
		{input: "HELP", expect: testCommand(t, "help")},
		{input: "hElP", expect: testCommand(t, "help")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := findCommand(tt.input)
			if !errors.Is(err, tt.err) {
				t.Fatalf(
					"findCommand(%q) Unexpected Error: %v",
					tt.input, err,
				)
			}
			if actual.name != tt.expect.name {
				t.Errorf(
					"findCommand(%q) Unexpected value: %v",
					tt.input, actual,
				)
			}
		})
	}
}

func TestCommandHandlers(t *testing.T) {
	tests := []struct {
		name string
		cmd string
		args []string
		pass func (*testUiSpy) (string, bool)
	} {
		{
			name: "help",
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleHelpCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "help-arg",
			cmd: "help",
			args: []string{"arg"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleCmdHelpCalled.called
				cw := tu.handleCmdHelpCalled.calledWith
				if !ca { return "handle not called", false }
				if cw != "arg" { return "handle called with wrong args", false }
				return "", true
			},
		},
		{
			name: "quit",
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleQuitCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "inventory",
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleInventoryCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "look",
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleLookCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "move",
			args: []string{"east"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "examine",
			args: []string{"obj-one"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "use",
			args: []string{"obj"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "use-item-on-obj",
			cmd: "use",
			args: []string{"item", "on", "obj"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "take",
			args: []string{"obj"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "drop",
			args: []string{"obj"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				if !ca { return "handle not called", false }
				return "", true
			},
		},
		{
			name: "move with bad object name should call handle target error",
			cmd: "move",
			args: []string{"obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "examine with bad object name should call handle target error",
			cmd: "examine",
			args: []string{"obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "use with bad object name should call handle target error",
			cmd: "use",
			args: []string{"obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "use-on with bad item name should call handle target error",
			cmd: "use",
			args: []string{"obj1", "on", "obj-one"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "use-on with bad object name should call handle target error",
			cmd: "use",
			args: []string{"obj-one", "on", "obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "take with bad object name should call handle target error",
			cmd: "take",
			args: []string{"obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
		{
			name: "drop with bad object name should call handle target error",
			cmd: "drop",
			args: []string{"obj1"},
			pass: func(tu *testUiSpy) (string, bool) {
				ca := tu.handleActionCalled.called
				eca := tu.handleErrorCalled.called
				ecw := tu.handleErrorCalled.calledWith
				if ca { return "handle action called", false }
				if !eca { return "handle target error not called", false }
				if !errors.Is(ecw, game.ErrInvalidIdentity) {
					return fmt.Sprintf("unexpected error: %v", ecw), false
				}
				return "", true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui :=  testUiSpy{}
			name := tt.name
			if tt.cmd != "" { name = tt.cmd }
			cmd := testCommand(t, name)
			cmd.handle(&ui, tt.args)
			if msg, ok := tt.pass(&ui); !ok {
				t.Fatalf("command(%q).handle(%v): %s", name, tt.args, msg)
			}
		})
	}
}
