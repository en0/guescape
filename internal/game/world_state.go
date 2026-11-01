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
	currentModeId 	Identity
	zones   		[]Zone
	objects 		[]Object
	flags			[]Identity
	modes 			[]Mode
	messages		map[MsgKey]string
}

// NewWorldState returns a validated WorldState, or an error if it is invalid.
func NewWorldState(
	startZoneId Identity,
	description string,
	flags		[]Identity,
	zones 		[]Zone,
	objects 	[]Object,
	modes 		[]Mode,
	messages	map[MsgKey]string,
) (WorldState, error) {
	ws := WorldState{
		currentZoneId: 	startZoneId,
		currentModeId:  Identity{},
		desc: 			description,
		flags: 			flags,
		zones: 			zones,
		objects: 		objects,
		modes: 			modes,
		messages:       messages,
	}
	if err := validateWorldState(ws); err != nil {
		return WorldState{}, err
	}
	return ws, nil
}

// Returns the world description
func (ws WorldState) Description() string {
	return ws.desc
}

// Returns the prompt state of the world
func (ws WorldState) Prompt() PromptView {
	var prompt string
	modeActive := ws.currentModeId != Identity{}
	if p, ok := ws.mode(ws.currentModeId); ok && modeActive {
		prompt = p.prompt
	} else {
		prompt = fmt.Sprintf("%s> ", ws.currentZoneId.String())
	}
	return PromptView{ Text: prompt, ModeActive: modeActive }
}

// Returns details about the current zone
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

// Returns a list of object ids currently in the player's inventory
func (ws WorldState) Inventory() []Identity {
	var result []Identity
	for _, o := range filter(ws.objects, func(o Object) bool {
		return o.inInventory()
	}) {
		result = append(result, o.id)
	}
	return result
}

// Check if the palyer has the given flag
func (ws WorldState) HasFlag(id Identity) bool {
	return exists(ws.flags, func(f Identity) bool {
		return f == id
	})
}

// Retrieve, from the given object, the value of the given property.
func (ws WorldState) ObjectProperty(objId Identity, propertyId Identity) (string, bool) {
	obj, found := ws.object(objId)
	if !found {
		return "", false
	}
	return obj.property(propertyId)
}

// Retrieve, from the given state value from the current mode.
func (ws WorldState) CurrentModeState(id Identity) (string, bool) {
	mode, found := ws.currentMode()
	if !found {
		return "", false
	}
	if v, ok := mode.states[id]; ok {
		return v, true
	}
	return "", false
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

func (ws WorldState) currentMode() (Mode, bool) {
	modeId := ws.currentModeId
	if modeId == (Identity{}) {
		return Mode{}, false
	}
	return ws.mode(modeId)
}

func (ws WorldState) mode(id Identity) (Mode, bool){
	return find(ws.modes, func(m Mode) bool {
		return m.id == id
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

func (ws WorldState) hasMode(id Identity) bool {
	return exists(ws.modes, func(m Mode) bool {
		return m.id == id
	})
}
