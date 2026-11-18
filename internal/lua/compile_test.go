package lua

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/en0/guescape/internal/game"
)

func TestCompileLua(t *testing.T) {
	tests := []struct {
		name string
		lua string
		expectErr error
	} {
		{
			name: "compile-lua-doesnt-error",
			lua: "local g = \"hello, world\"",
		},
		{
			name: "compile-lua-doesnt-error",
			lua: "hello, world",
			expectErr: ErrCompileFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompileLua(tt.lua)
			if !errors.Is(err, tt.expectErr) {
				t.Fatalf("CompileLua(...) err= %v, wants %v", err, tt.expectErr)
			}
		})
	}
}

func TestRunLua(t *testing.T) {
	tests := []struct {
		name string
		lua string
		args []string
		expectErr error
		ctxSpy *MockScriptContext
		expectEffects []game.Effect
	} {
		{
			name: "run-completes",
			lua: "local g = \"hello, world\"",
		},
		{
			name: "has-flag-calls-ctx",
			lua: `hasFlag("flag-one")`,
			ctxSpy: &MockScriptContext{
				hasFlag: true,
				hasFlagCalled: true,
				hasFlagCalledWith: testIdentity(t, "flag-one"),
			},
		},
		{
			name: "get-flag-err-causes-abort",
			lua: `hasFlag("flag1")`,
			expectErr: ErrExecuteScript,
		},
		{
			name: "get-mode-state-calls-ctx",
			lua: `getModeState("mode-one", "state-one")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
				stateString: "state-one-value",
				stateBool: true,
				stateCalled: true,
				stateCalledWith: testIdentities(t, "mode-one", "state-one"),
			},
		},
		{
			name: "get-mode-state-err-causes-abort",
			lua: `getModeState("mode-one", "state1")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
			},
			expectErr: ErrExecuteScript,
		},
		{
			name: "get-object-state-calls-ctx",
			lua: `getObjectState("obj-one", "state-one")`,
			ctxSpy: &MockScriptContext{
				objectStateString: "state-one-value",
				objectStateBool: true,
				objectStateCalled: true,
				objectStateCalledWith: testIdentities(t, "obj-one", "state-one"),
			},
		},
		{
			name: "get-object-state-obj-err-causes-abort",
			lua: `getObjectState("obj1", "state-one")`,
			expectErr: ErrExecuteScript,
		},
		{
			name: "get-object-state-err-causes-abort",
			lua: `getObjectState("obj-one", "state1")`,
			expectErr: ErrExecuteScript,
		},
		{
			name: "open-table-works",
			lua: `
				local t = {1, 2, 3}
				table.insert(t, 4)
				assert(#t == 4)
			`,
		},
		{
			name: "open-string-works",
			lua: `assert(string.upper("abc") == "ABC")`,
		},
		{
			name: "open-math-works",
			lua: `assert(math.floor(math.pi) == 3)`,
		},
		{
			name: "print-fails",
			expectErr: ErrExecuteScript,
			lua: `print("bad")`,
		},
		{
			name: "load-fails",
			expectErr: ErrExecuteScript,
			lua: `
				load(function()
					if done then return nil end
					done = true
					return "return 1 + 2"
				end)
			`,
		},
		{
			name: "loadstring-fails",
			expectErr: ErrExecuteScript,
			lua: `loadstring("return 1 + 2")`,
		},
		{
			name: "dofile-fails",
			expectErr: ErrExecuteScript,
			lua: `dofile("/dev/null")`,
		},
		{
			name: "say-effect",
			lua: `say("say-one")`,
			expectEffects: []game.Effect{
				game.SayEffect{Text: "say-one"},
			},
		},
		{
			name: "set-current-zone-effect",
			lua: `setCurrentZone("zone-one")`,
			expectEffects: []game.Effect{
				game.CurrentZoneEffect{ZoneId: testIdentity(t, "zone-one")},
			},
		},
		{
			name: "add-flag-effect",
			lua: `addFlag("flag-one")`,
			expectEffects: []game.Effect{
				game.AddFlagEffect{Flag: testIdentity(t, "flag-one")},
			},
		},
		{
			name: "remove-flag-effect",
			lua: `removeFlag("flag-one")`,
			expectEffects: []game.Effect{
				game.RemoveFlagEffect{Flag: testIdentity(t, "flag-one")},
			},
		},
		{
			name: "set-object-state-effect",
			lua: `setObjectState("obj-one", "state-one", "state-one-value")`,
			expectEffects: []game.Effect{
				game.SetObjectStateEffect{
					ObjectId: testIdentity(t, "obj-one"),
					StateId: testIdentity(t, "state-one"),
					Value: "state-one-value",
				},
			},
		},
		{
			name: "set-visibility-effect",
			lua: `setVisibility("obj-one", true)`,
			expectEffects: []game.Effect{
				game.ObjectVisibleEffect{
					ObjectId: testIdentity(t, "obj-one"),
					Visible: true,
				},
			},
		},
		{
			name: "set-targetable-effect",
			lua: `setTargetable("obj-one", true)`,
			expectEffects: []game.Effect{
				game.ObjectTargetableEffect{
					ObjectId: testIdentity(t, "obj-one"),
					Targetable: true,
				},
			},
		},
		{
			name: "set-carriable-effect",
			lua: `setCarriable("obj-one", true)`,
			expectEffects: []game.Effect{
				game.ObjectCarriableEffect{
					ObjectId: testIdentity(t, "obj-one"),
					Carriable: true,
				},
			},
		},
		{
			name: "set-object-location-effect",
			lua: `setObjectLocation("obj-one", "zone-one")`,
			expectEffects: []game.Effect{
				game.ObjectLocationEffect{
					ObjectId: testIdentity(t, "obj-one"),
					ZoneId: testIdentity(t, "zone-one"),
				},
			},
		},
		{
			name: "end-game-effect",
			lua: `endGame(true, "unit test passed")`,
			expectEffects: []game.Effect{
				game.EndGameEffect{
					IsWin: true,
					Text: "unit test passed",
				},
			},
		},
		{
			name: "set-mode-effect",
			lua: `activateMode("mode-one")`,
			expectEffects: []game.Effect{
				game.ActivateModeEffect{
					ModeId: testIdentity(t, "mode-one"),
				},
			},
		},
		{
			name: "clear-mode-effect",
			lua: `clearMode()`,
			expectEffects: []game.Effect{
				game.ClearModeEffect{},
			},
		},
		{
			name: "set-mode-state-effect",
			lua: `setModeState("mode-one", "state-one", "state-one-value")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
			},
			expectEffects: []game.Effect{
				game.SetModeStateEffect{
					ModeId: testIdentity(t, "mode-one"),
					State: testIdentity(t, "state-one"),
					Value: "state-one-value",
				},
			},
		},
		{
			name: "set-prompt-effect",
			lua: `setPrompt("mode-one-value")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
			},
			expectEffects: []game.Effect{
				game.SetModePromptEffect{
					ModeId: testIdentity(t, "mode-one"),
					Value: "mode-one-value",
				},
			},
		},
		{
			name: "set-prompt-effect",
			args: []string{"hello,", "lua!"},
			lua: "for i, a in ipairs(args) do say(a) end",
			expectEffects: []game.Effect{
				game.SayEffect{Text: "hello,"},
				game.SayEffect{Text: "lua!"},
			},
		},
		{
			name: "get-state-calls-ctx-shorthand",
			lua: `get("state-one")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
				stateString: "state-one-value",
				stateBool: true,
				stateCalled: true,
				stateCalledWith: testIdentities(t, "mode-one", "state-one"),
			},
		},
		{
			name: "set-mode-state-effect-shorthand",
			lua: `set("state-one", "state-one-value")`,
			ctxSpy: &MockScriptContext{
				currentModeId: testIdentity(t, "mode-one"),
			},
			expectEffects: []game.Effect{
				game.SetModeStateEffect{
					ModeId: testIdentity(t, "mode-one"),
					State: testIdentity(t, "state-one"),
					Value: "state-one-value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := CompileLua(tt.lua)
			if err != nil {
				t.Fatalf("CompileLua(...) err = %v", err)
			}
			ctx := testMockScriptContextFromContext(tt.ctxSpy)
			actual, err := s.Run(ctx, tt.args)
			if !errors.Is(err, tt.expectErr) {
				t.Fatalf(
					"s.Run(%#v) err= %v, wants %v",
					tt.args, err, tt.expectErr)
			}
			if tt.ctxSpy != nil && !reflect.DeepEqual(ctx, tt.ctxSpy) {
				t.Fatalf(
					"ctx: got %v, wanted %v",
					ctx, tt.ctxSpy)
			}
			if !slices.Equal(actual, tt.expectEffects) {
				t.Fatalf(
					"s.Run(%#v) eff= %v, wants %v",
					tt.args, actual, tt.expectEffects)
			}

		})
	}
}

