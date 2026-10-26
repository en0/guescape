package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

// objectBlock

type objectBlock struct {
	Id 			string			`hcl:"id,label"`
	ZoneId 		string			`hcl:"zone_id"`
	Desc 		*string			`hcl:"desc"`
	Visible 	*bool			`hcl:"visible"`
	Targetable 	*bool			`hcl:"targetable"`
	Carriable 	*bool			`hcl:"carriable"`
	Properties 	[]propertyBlock	`hcl:"property,block"`
	Actions 	[]actionBlock	`hcl:"action,block"`
}

type propertyBlock struct {
	Key 	string	`hcl:"id,label"`
	Value 	string	`hcl:"value"`
}

type actionBlock struct {
	Verb 	string		`hcl:"verb,label"`
	Rules 	[]ruleBlock	`hcl:"rule,block"`
	// Mutually exclusive short-hand
	//Effects []Effect
	//EffectsBlock []EffectBlock	`hcl:"effect,block"`
}

func (o objectBlock) toDomain() (game.Object, error) {
	var errs []error

	oid, err := game.NewIdentity(o.Id)
	if err != nil {errs = append(errs, err)}

	zid, err := game.NewIdentity(o.ZoneId)
	if err != nil {errs = append(errs, err)}

	var desc string
	if o.Desc != nil {desc = *o.Desc}

	var pm map[game.Identity]string
	for _, p := range o.Properties {
		v, err := game.NewIdentity(p.Key)
		if err != nil {errs = append(errs, err); continue}
		if pm == nil {pm = make(map[game.Identity]string)}
		pm[v] = p.Value
	}

	var am map[game.ObjectVerb][]game.Rule
	for _, a := range o.Actions {
		verb, err := NewObjectVerb(a.Verb)
		if err != nil {errs = append(errs, err); continue}
		if am == nil {
			am = make(map[game.ObjectVerb][]game.Rule)
		}
		for _, rb := range a.Rules {
			v, err := rb.toDomain()
			if err != nil {errs = append(errs, err); continue}
			am[verb] = append(am[verb], v)
		}
	}
	return game.NewObject(
		oid, zid, desc,
		o.Visible, o.Targetable, o.Carriable,
		pm, am,
	), errors.Join(errs...)
}
