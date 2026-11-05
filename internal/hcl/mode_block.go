package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

type modeBlock struct {
	Id 		string				`hcl:"id,label"`
	Prompt	*string				`hcl:"prompt"`
	Lua 	*string				`hcl:"lua"`
	States  *map[string]string 	`hcl:"state"`
	luaScript game.CompiledScript
}

func (m modeBlock) toDomain() (game.Mode, error) {
	var errs []error
	id, err := game.NewIdentity(m.Id)
	if err != nil {
		errs = append(errs, err)
	}
	var states map[game.Identity]string
	for k, v := range *m.States {
		s, err := game.NewIdentity(k)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if states == nil {
			states = make(map[game.Identity]string)
		}
		states[s] = v
	}
	var prompt string
	if m.Prompt != nil {
		prompt = *m.Prompt
	}
	return game.NewMode(id, prompt, states, m.luaScript), errors.Join(errs...)
}
