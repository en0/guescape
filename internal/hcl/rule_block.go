package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

// ruleBlock

type ruleBlock struct {
	ConditionGroups []conditionGroupBlock	`hcl:"when,block"`
	EffectsBlock 	[]effectBlock			`hcl:"effect,block"`
	Effects 		[]effect
}

func (r ruleBlock) toDomain() (game.Rule, error) {
	var errs []error

	var cg []game.ConditionGroup
	for _, c := range r.ConditionGroups {
		v, err := c.toDomain()
		if err != nil {errs = append(errs, err); continue}
		cg = append(cg, v)
	}

	var ef []game.Effect
	for _, e := range r.Effects {
		v, err := e.toDomain()
		if err != nil {errs = append(errs, err); continue}
		ef = append(ef, v)
	}

	return game.NewRule(cg, ef), errors.Join(errs...)
}

// conditionGroupBlock

type conditionGroupBlock struct {
	PropertyConditions 	[]propertyConditionBlock	`hcl:"property,block"`
	FlagConditions 		[]flagConditionBlock		`hcl:"flag,block"`
	UseWith 			*string						`hcl:"use_with"`
}

func (g conditionGroupBlock) toDomain() (game.ConditionGroup, error) {
	var errs []error

	var pc []game.PropertyCondition
	for _, c := range g.PropertyConditions {
		v, err := c.toDomain()
		if err != nil {errs = append(errs, err); continue}
		pc = append(pc, v)
	}

	var fc []game.FlagCondition
	for _, c := range g.FlagConditions {
		v, err := c.toDomain()
		if err != nil {errs = append(errs, err); continue}
		fc = append(fc, v)
	}

	var useWith game.Identity
	if g.UseWith != nil {
		var err error
		useWith, err = game.NewIdentity(*g.UseWith)
		if err != nil {errs = append(errs, err)}
	}

	return game.NewConditionGroup(useWith, fc, pc), errors.Join(errs...)
}

// propertyConditionBlock

type propertyConditionBlock struct {
	Id 		string	`hcl:"id,label"`
	Value 	string	`hcl:"value"`
}

func (p propertyConditionBlock) toDomain() (game.PropertyCondition, error) {
	id, err := game.NewIdentity(p.Id)
	return game.NewPropertyCondition(id, p.Value), err
}

// flagConditionBlock

type flagConditionBlock struct {
	Id 		string	`hcl:"id,label"`
	Value 	bool	`hcl:"value"`
}

func (f flagConditionBlock) toDomain() (game.FlagCondition, error) {
	id, err := game.NewIdentity(f.Id)
	return game.NewFlagCondition(id, f.Value), err
}
