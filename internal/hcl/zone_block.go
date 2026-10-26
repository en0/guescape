package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

// zoneBlock

type zoneBlock struct {
	Id 		string  	`hcl:"id,label"`
	Desc 	string		`hcl:"desc"`
	Exits 	[]exitBlock	`hcl:"exit,block"`
}

func (z zoneBlock) toDomain() (game.Zone, error) {
	var errs []error

	id, err := game.NewIdentity(z.Id)
	if err != nil {errs = append(errs, err)}

	var em map[game.Identity]game.ZoneExit
	for _, eb := range z.Exits {
		id, err := game.NewIdentity(eb.Id)
		if err != nil {errs = append(errs, err); continue}
		v, err := eb.toDomain()
		if err != nil {errs = append(errs, err); continue}
		if em == nil {em = make(map[game.Identity]game.ZoneExit)}
		em[id] = v
	}

	return game.NewZone(id, z.Desc, em), errors.Join(errs...)
}

// exitBlock

type exitBlock struct {
	Id 		string 		`hcl:"id,label"`
	ZoneId 	string 		`hcl:"to"`
	Flags 	[]string 	`hcl:"flags"`
}

func (e exitBlock) toDomain() (game.ZoneExit, error) {
	var errs []error

	id, err := game.NewIdentity(e.ZoneId)
	if err != nil {errs = append(errs, err)}

	var flags []game.Identity
	for _, f := range e.Flags {
		v, err := game.NewIdentity(f)
		if err != nil {errs = append(errs, err); continue}
		flags = append(flags, v)
	}

	return game.NewZoneExit(id, flags), errors.Join(errs...)
}
