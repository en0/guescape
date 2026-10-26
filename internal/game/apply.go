package game

import (
	"errors"
	"maps"
	"slices"
)

var errEffectObjectNotFound = errors.New("effect references an unknown object")
var errEffectPropertyNotFound = errors.New("effect references an unknown property")

// Apply returns a new WorldState with ef applied to ws, leaving ws unchanged.
// SayEffect and EndGameEffect are no-ops. On error, discard the returned state.
func Apply(ws WorldState, ef Effect) (WorldState, error){
	var err error
	switch e := ef.(type) {
	case CurrentZoneEffect:
		ws.currentZoneId = e.ZoneId
	case AddFlagEffect:
		ws = applyAddFlagEffect(ws, e)
	case RemoveFlagEffect:
		ws = applyRemoveFlagEffect(ws, e)
	case ObjectPropertyEffect:
		ws, err = applyObjectPropertyEffect(ws, e)
	case ObjectVisibleEffect:
		ws, err = applyObjectVisibleEffect(ws, e)
	case ObjectTargetableEffect:
		ws, err = applyObjectTargetableEffect(ws, e)
	case ObjectLocationEffect:
		ws, err = applyObjectLocationEffect(ws, e)
	}
	if err != nil {return WorldState{}, err}
	return ws, nil
}

func applyAddFlagEffect(ws WorldState, e AddFlagEffect) WorldState {
	if !ws.hasFlag(e.Flag) {
		ws.flags = append(ws.flags, e.Flag)
	}
	return ws
}

func applyRemoveFlagEffect(ws WorldState, e RemoveFlagEffect) WorldState {
	ws.flags = slices.DeleteFunc(slices.Clone(ws.flags), func(f Identity) bool {
		return f == e.Flag
	})
	return ws
}

func applyObjectPropertyEffect(ws WorldState, e ObjectPropertyEffect) (WorldState, error) {
	i := slices.IndexFunc(ws.objects, func(o Object) bool {return o.id == e.ObjectId})
	if i > -1 {
		objects := slices.Clone(ws.objects)
		if objects[i].hasProperty(e.PropertyId) {
			objects[i].properties = maps.Clone(objects[i].properties)
			objects[i].properties[e.PropertyId] = e.Value
			ws.objects = objects
			return ws, nil
		}
		return WorldState{}, errEffectPropertyNotFound
	}
	return WorldState{}, errEffectObjectNotFound
}

func applyObjectVisibleEffect(ws WorldState, e ObjectVisibleEffect) (WorldState, error) {
	if updateObject(&ws, e.ObjectId, func(o *Object) {
		o.visible = e.Visible
	}) {
		return ws, nil
	}
	return WorldState{}, errEffectObjectNotFound
}

func applyObjectTargetableEffect(ws WorldState, e ObjectTargetableEffect) (WorldState, error) {
	if updateObject(&ws, e.ObjectId, func(o *Object) {
		o.targetable = e.Targetable
	}) {
		return ws, nil
	}
	return WorldState{}, errEffectObjectNotFound
}

func applyObjectLocationEffect(ws WorldState, e ObjectLocationEffect) (WorldState, error) {
	if updateObject(&ws, e.ObjectId, func(o *Object) {
		o.zoneId = e.ZoneId
	}) {
		return ws, nil
	}
	return WorldState{}, errEffectObjectNotFound
}

// Updates an object using a clone-on-write strategy to avoid mutating the WorldState.
func updateObject(ws *WorldState, id Identity, edit func(*Object)) bool {
	i := slices.IndexFunc(ws.objects, func(o Object) bool {return o.id == id})
	if i < 0 { return false }
	ws.objects = slices.Clone(ws.objects)
	edit(&ws.objects[i])
	return true
}
