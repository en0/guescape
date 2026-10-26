package hcl

import (
	"errors"
	"fmt"
	"os"

	"github.com/en0/guescape/internal/game"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type worldRoot struct {
	Game worldBlock `hcl:"game,block"`
}

// LoadFile loads the world state from the local file system.
// The format is expected to by hcl.v2
func LoadFile(filename string) (game.WorldState, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return game.WorldState{}, err
	}

	body, err := parseHCL(src, filename)
	if err != nil {
		return game.WorldState{}, err
	}

	world, err := decodeWorld(body)
	if err != nil {
		return game.WorldState{}, err
	}

	state, err := world.toDomain()
	if err != nil {
		return game.WorldState{}, err
	}

	return state, nil
}

func parseHCL(src []byte, filename string) (hcl.Body, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(src, filename)
	if diags.HasErrors() {
		errs := appendDiagsAs(nil, diags, parsingErrorMapper)
		return nil, errors.Join(errs...)
	}
	return file.Body, nil
}

func decodeWorld(body hcl.Body) (worldBlock, error) {
	var root worldRoot
	diags := gohcl.DecodeBody(body, nil, &root)
	if diags.HasErrors() {
		errs := appendDiagsAs(nil, diags, decodingErrorMapper)
		return worldBlock{}, errors.Join(errs...)
	}
	return root.Game, resolveEffects(&root.Game)
}

func resolveEffects(g *worldBlock) error {
	var errs []error
	for i := range g.Objects {
		object := &g.Objects[i]
		for j := range object.Actions {
			action := &object.Actions[j]
			for k := range action.Rules {
				rule := &action.Rules[k]
				for l := range rule.EffectsBlock {
					effect, diags := loadEffect(rule.EffectsBlock[l])
					if diags.HasErrors() {
						errs = appendDiagsAs(errs, diags, decodingErrorMapper)
						continue
					}
					rule.Effects = append(rule.Effects, effect)
				}
				rule.EffectsBlock = nil
			}
		}
	}
	return errors.Join(errs...)
}

func loadEffect(b effectBlock) (effect, hcl.Diagnostics) {
	switch b.Kind {
	case "say":
		var v sayEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "property":
		var v objectPropertyEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "object_targetable":
		var v objectTargetableEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "add_flag":
		var v addFlagEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "remove_flag":
		var v removeFlagEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "end_game":
		var v endGameEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "current_zone":
		var v currentZoneEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "object_visible":
		var v objectVisibleEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	case "object_zone":
		var v objectLocationEffectBlock
		d := gohcl.DecodeBody(b.Body, nil, &v)
		return v, d
	default:
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "unknown effect type",
			Detail:   fmt.Sprintf("%q is not a recognized effect discriminator", b.Kind),
		}}
	}
}
