package game

import (
	"fmt"
	"maps"
	"slices"
)

// WorldState holds the current state of the game world.
type WorldState struct {
	desc	 		string
	currentZoneId 	Identity
	zones   		[]Zone
	objects 		[]Object
	flags			[]Identity
	messages		map[MsgKey]string
}

// NewWorldState returns a validated WorldState, or an error if it is invalid.
func NewWorldState(
	startZoneId Identity,
	description string,
	flags		[]Identity,
	zones 		[]Zone,
	objects 	[]Object,
	messages	map[MsgKey]string,
) (WorldState, error) {
	ws := WorldState{
		currentZoneId: 	startZoneId,
		desc: 			description,
		flags: 			flags,
		zones: 			zones,
		objects: 		objects,
		messages:       messages,
	}
	if err := validateWorldState(ws); err != nil {
		return WorldState{}, err
	}
	return ws, nil
}

func (ws WorldState) Description() string {
	return ws.desc
}

func (ws WorldState) ViewCurrentZone() ZoneView {

	var objects []Identity
	for _, o := range filter(ws.objects, func(o Object) bool {
		return o.visible && o.zoneId == ws.currentZoneId
	}) {
		objects = append(objects, o.id)
	}

	zone := ws.currentZone()
	exits := slices.Collect(maps.Keys(zone.exits))

	return ZoneView{
		Id: zone.id,
		Desc: zone.desc,
		Exits: exits,
		Objects: objects,
	}
}

func (ws WorldState) Inventory() []Identity {
	var result []Identity
	for _, o := range filter(ws.objects, func(o Object) bool {
		return o.inInventory()
	}) {
		result = append(result, o.id)
	}
	return result
}

func (ws WorldState) message(key MsgKey) string {
	if val, ok := ws.messages[key]; ok { return val }
	return message(key)
}

func (ws WorldState) currentZone() (Zone) {
	zone, ok := ws.zone(ws.currentZoneId)
	if ok { return zone }
	panic(fmt.Errorf("State.FindCurrentZone: state has no current zone"))
}

func (ws WorldState) zone(id Identity) (Zone, bool) {
	return find(ws.zones, func(z Zone) bool {
		return z.id == id
	})
}

func (ws WorldState) object(id Identity) (Object, bool) {
	return find(ws.objects, func(o Object) bool {
		return o.id == id
	})
}

func (ws WorldState) actionableObject(id Identity) (Object, bool) {
	return find(ws.objects, func(o Object) bool {
		if o.id != id { return false }
		inCurrentZone := o.zoneId == ws.currentZoneId && o.targetable
		return inCurrentZone || o.inInventory()
	})
}

func (ws WorldState) hasZone(id Identity) bool {
	return exists(ws.zones, func(z Zone) bool {
		return z.id == id
	})
}

func (ws WorldState) hasObject(id Identity) bool {
	return exists(ws.objects, func(z Object) bool {
		return z.id == id
	})
}

func (ws WorldState) hasFlag(id Identity) bool {
	return exists(ws.flags, func(f Identity) bool {
		return f == id
	})
}
