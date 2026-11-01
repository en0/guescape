package lua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/en0/guescape/internal/game"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var ErrParseFailed = errors.New("could not compile lua")
var ErrCompileFailed = errors.New("could not compile lua")
var ErrExecuteScript = errors.New("failed to execute lua")

func CompileLua(src string) (game.CompiledScript, error) {
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()
	chunk, err := parse.Parse(strings.NewReader(src), "mode-script")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCompileFailed, err)
	}
	proto, err := lua.Compile(chunk, "mode-script")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCompileFailed, err)
	}
	return &CompiledScript{
		proto: proto,
	}, nil
}

type CompiledScript struct {
	proto *lua.FunctionProto
}

func (s *CompiledScript) Run(ctx game.ScriptContext, args []string) ([]game.Effect, error) {
	// Remove all built-ins from VM
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	// Selectivly re-enable core lib elements
	for _, f := range []struct {
		name string
		fn lua.LGFunction
	} {
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.StringLibName, lua.OpenString},
		{lua.MathLibName, lua.OpenMath},
	} {
		L.Push(L.NewFunction(f.fn))
		L.Push(lua.LString(f.name))
		L.Call(1, 0)
	}

	// Explicitly disable dangerous commands
	for _, name := range []string{
		"print",
		"load",
		"loadstring",
		"dofile",
		"collectgarbage",
	} { L.SetGlobal(name, lua.LNil) }

	// Register context globals
	L.SetGlobal("hasFlag", L.NewFunction(lgHasFlag(ctx)))
	L.SetGlobal("getState", L.NewFunction(lgGetState(ctx)))
	L.SetGlobal("getProperty", L.NewFunction(lgGetProperty(ctx)))

	effects := []game.Effect{}

	// Register effect handlers
	L.SetGlobal("say", L.NewFunction(lgSay(&effects)))
	L.SetGlobal("setCurrentZone", L.NewFunction(lgSetCurrentZone(&effects)))
	L.SetGlobal("addFlag", L.NewFunction(lgAddFlag(&effects)))
	L.SetGlobal("removeFlag", L.NewFunction(lgRemoveFlag(&effects)))
	L.SetGlobal("setProperty", L.NewFunction(lgSetProperty(&effects)))
	L.SetGlobal("setVisibility", L.NewFunction(lgSetVisibility(&effects)))
	L.SetGlobal("setTargetable", L.NewFunction(lgSetTargetable(&effects)))
	//L.SetGlobal("setCarriable", L.NewFunction(lgSetCarriable(&effects)))
	L.SetGlobal("setObjectLocation", L.NewFunction(lgSetObjectLocation(&effects)))
	L.SetGlobal("endGame", L.NewFunction(lgEndGame(&effects)))
	L.SetGlobal("setMode", L.NewFunction(lgSetMode(&effects)))
	L.SetGlobal("clearMode", L.NewFunction(lgClearMode(&effects)))
	L.SetGlobal("setState", L.NewFunction(lgSetState(&effects)))
	L.SetGlobal("setPrompt", L.NewFunction(lgSetPrompt(&effects)))

	L.Push(L.NewFunctionFromProto(s.proto))
	err := L.PCall(0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrExecuteScript, err)
	}

	return effects, nil
}

func lgHasFlag(c game.ScriptContext) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("hasFlag: invalid identity %q: %v", lv, err)
			return 0
		}
		result := c.HasFlag(id)
		l.Push(lua.LBool(result))
		return 1
	}
}

func lgGetState(c game.ScriptContext) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("getState: invalid identity %q: %v", lv, err)
			return 0
		}
		if v, ok := c.CurrentModeState(id); ok {
			l.Push(lua.LString(v))
			return 1
		}
		return 0
	}
}

func lgGetProperty(c game.ScriptContext) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		oid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("getState: invalid identity %q: %v", lv1, err)
			return 0
		}
		pid, err := game.NewIdentity(lv2)
		if err != nil {
			l.RaiseError("getState: invalid identity %q: %v", lv2, err)
			return 0
		}
		if v, ok := c.ObjectProperty(oid, pid); ok {
			l.Push(lua.LString(v))
			return 1
		}
		return 0
	}
}

func lgSay(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		e := game.SayEffect{Text: lv}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetCurrentZone(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("setCurrentZone: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.CurrentZoneEffect{ZoneId: id}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgAddFlag(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("addFlag: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.AddFlagEffect{Flag: id}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgRemoveFlag(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("removeFlag: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.RemoveFlagEffect{Flag: id}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetProperty(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		lv3 := l.ToString(3)
		oid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setProperty: invalid identity %q: %v", lv1, err)
			return 0
		}
		pid, err := game.NewIdentity(lv2)
		if err != nil {
			l.RaiseError("setProperty: invalid identity %q: %v", lv2, err)
			return 0
		}
		e := game.ObjectPropertyEffect{
			ObjectId: oid,
			PropertyId: pid,
			Value: lv3,
		}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetVisibility(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToBool(2)
		id, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setVisibility: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.ObjectVisibleEffect{ObjectId: id, Visible: lv2}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetTargetable(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToBool(2)
		id, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setTargetable: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.ObjectTargetableEffect{ObjectId: id, Targetable: lv2}
		*effecs = append(*effecs, e)
		return 0
	}
}

// func lgSetCarriable(effecs *[]game.Effect) lua.LGFunction {
// 	return func(l *lua.LState) int {
// 		lv1 := l.ToString(1)
// 		lv2 := l.ToBool(2)
// 		id, err := game.NewIdentity(lv1)
// 		if err != nil {
// 			l.RaiseError("setCarriable: invalid identity %q: %v", lv1, err)
// 			return 0
// 		}
// 		e := game.ObjectCarriableEffect{ObjectId: id, Carriable: lv2}
// 		*effecs = append(*effecs, e)
// 		return 0
// 	}
// }

func lgSetObjectLocation(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		oid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setObjectLocation: invalid identity %q: %v", lv1, err)
			return 0
		}
		zid, err := game.NewIdentity(lv2)
		if err != nil {
			l.RaiseError("setObjectLocation: invalid identity %q: %v", lv2, err)
			return 0
		}
		e := game.ObjectLocationEffect{ObjectId: oid, ZoneId: zid}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgEndGame(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToBool(1)
		lv2 := l.ToString(2)
		e := game.EndGameEffect{IsWin: lv1, Text: lv2}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetMode(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("setMode: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.SetModeEffect{ModeId: id}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgClearMode(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		e := game.ClearModeEffect{}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetState(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		lv3 := l.ToString(3)
		mid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setState: invalid identity %q: %v", lv1, err)
			return 0
		}
		sid, err := game.NewIdentity(lv2)
		if err != nil {
			l.RaiseError("setState: invalid identity %q: %v", lv2, err)
			return 0
		}
		e := game.SetModeStateEffect{ModeId: mid, State: sid, Value: lv3}
		*effecs = append(*effecs, e)
		return 0
	}
}

func lgSetPrompt(effecs *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		mid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setPrompt: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.SetModePromptEffect{ModeId: mid, Value: lv2}
		*effecs = append(*effecs, e)
		return 0
	}
}
