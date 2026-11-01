package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

type modeBlock struct {
	Id 		string			`hcl:"id,label"`
	Prompt	*string			`hcl:"prompt"`
	Lua 	*string			`hcl:"lua"`
	States 	[]stateBlock	`hcl:"state,block"`
	luaScript game.CompiledScript
}

type stateBlock struct {
	Id 		string	`hcl:"id,label"`
	Value	string	`hcl:"value"`
}

func (m modeBlock) toDomain() (game.Mode, error) {
	var errs []error
	id, err := game.NewIdentity(m.Id)
	if err != nil {
		errs = append(errs, err)
	}
	var states map[game.Identity]string
	for _, sb := range m.States {
		sid, err := game.NewIdentity(sb.Id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if states == nil {
			states = make(map[game.Identity]string)
		}
		states[sid] = sb.Value
	}
	var prompt string
	if m.Prompt != nil {
		prompt = *m.Prompt
	}
	return game.NewMode(id, prompt, states, m.luaScript), errors.Join(errs...)
}
