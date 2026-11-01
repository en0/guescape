package lua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/en0/guescape/internal/game"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var ErrCompileFailed = errors.New("could not compile lua")
var ErrExecuteScript = errors.New("failed to execute lua")

func CompileLua(src string) (game.CompiledScript, error) {
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
	L.SetGlobal("setCarriable", L.NewFunction(lgSetCarriable(&effects)))
	L.SetGlobal("setObjectLocation", L.NewFunction(lgSetObjectLocation(&effects)))
	L.SetGlobal("endGame", L.NewFunction(lgEndGame(&effects)))
	L.SetGlobal("setMode", L.NewFunction(lgSetMode(&effects)))
	L.SetGlobal("clearMode", L.NewFunction(lgClearMode(&effects)))
	L.SetGlobal("setState", L.NewFunction(lgSetState(ctx, &effects)))
	L.SetGlobal("setPrompt", L.NewFunction(lgSetPrompt(ctx, &effects)))

	// Register args as global table
	argTable := L.NewTable()
	for _, a := range args { argTable.Append(lua.LString(a)) }
	L.SetGlobal("args", argTable)

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
		if v, ok := c.ModeState(c.CurrentModeId(), id); ok {
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

func lgSay(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		e := game.SayEffect{Text: lv}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetCurrentZone(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("setCurrentZone: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.CurrentZoneEffect{ZoneId: id}
		*effects = append(*effects, e)
		return 0
	}
}

func lgAddFlag(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("addFlag: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.AddFlagEffect{Flag: id}
		*effects = append(*effects, e)
		return 0
	}
}

func lgRemoveFlag(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("removeFlag: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.RemoveFlagEffect{Flag: id}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetProperty(effects *[]game.Effect) lua.LGFunction {
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
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetVisibility(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToBool(2)
		id, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setVisibility: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.ObjectVisibleEffect{ObjectId: id, Visible: lv2}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetTargetable(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToBool(2)
		id, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setTargetable: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.ObjectTargetableEffect{ObjectId: id, Targetable: lv2}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetCarriable(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToBool(2)
		id, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setCarriable: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.ObjectCarriableEffect{ObjectId: id, Carriable: lv2}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetObjectLocation(effects *[]game.Effect) lua.LGFunction {
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
		*effects = append(*effects, e)
		return 0
	}
}

func lgEndGame(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToBool(1)
		lv2 := l.ToString(2)
		e := game.EndGameEffect{IsWin: lv1, Text: lv2}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetMode(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		id, err := game.NewIdentity(lv)
		if err != nil {
			l.RaiseError("setMode: invalid identity %q: %v", lv, err)
			return 0
		}
		e := game.SetModeEffect{ModeId: id}
		*effects = append(*effects, e)
		return 0
	}
}

func lgClearMode(effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		e := game.ClearModeEffect{}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetState(ctx game.ScriptContext, effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv1 := l.ToString(1)
		lv2 := l.ToString(2)
		sid, err := game.NewIdentity(lv1)
		if err != nil {
			l.RaiseError("setState: invalid identity %q: %v", lv1, err)
			return 0
		}
		e := game.SetModeStateEffect{ModeId: ctx.CurrentModeId(), State: sid, Value: lv2}
		*effects = append(*effects, e)
		return 0
	}
}

func lgSetPrompt(ctx game.ScriptContext, effects *[]game.Effect) lua.LGFunction {
	return func(l *lua.LState) int {
		lv := l.ToString(1)
		e := game.SetModePromptEffect{ModeId: ctx.CurrentModeId(), Value: lv}
		*effects = append(*effects, e)
		return 0
	}
}
