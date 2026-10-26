package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

type worldBlock struct {
	Desc 		string			`hcl:"desc"`
	StartZoneId string 			`hcl:"start_zone_id"`
	Zones 		[]zoneBlock		`hcl:"zone,block"`
	Objects 	[]objectBlock	`hcl:"object,block"`
	Messages	[]messageBlock	`hcl:"message,block"`
}

type messageBlock struct {
	Key string	`hcl:"id,label"`
	Text string	`hcl:"text"`
}

func (wb worldBlock) toDomain() (game.WorldState, error) {
	var errs []error

	szid, err := game.NewIdentity(wb.StartZoneId)
	if err != nil {errs = append(errs, err)}

	var zones []game.Zone
	for _, zb := range wb.Zones {
		v, err := zb.toDomain()
		if err != nil {errs = append(errs, err); continue}
		zones = append(zones, v)
	}

	var objects []game.Object
	for _, ob := range wb.Objects {
		v, err := ob.toDomain()
		if err != nil {errs = append(errs, err); continue}
		objects = append(objects, v)
	}

	var messages map[game.MsgKey]string
	for _, mb := range wb.Messages {
		key, err := NewMessageKey(mb.Key)
		if err != nil {errs = append(errs, err); continue}
		if messages == nil {messages = make(map[game.MsgKey]string)}
		messages[key] = mb.Text
	}

	if len(errs) > 0 {
		return game.WorldState{}, errors.Join(errs...)
	}

	return game.NewWorldState(szid, wb.Desc, nil, zones, objects, messages)
}
