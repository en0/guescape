package game

import (
	"errors"
)

var errObjectZoneNotFound = errors.New("object references an unknown zone")
var errCurrentZoneEffectZoneNotFound = errors.New("current zone effect references an unknown zone")
var errObjectLocationEffectZoneNotFound = errors.New("object location effect references an unknown zone")
var errObjectLocationEffectObjectNotFound = errors.New("object location effect references an unknown object")
var errObjectTargetableEffectObjectNotFound = errors.New("object targetable effect references and unknown object")
var errObjectVisibleEffectObjectNotFound = errors.New("object visible effect references an unknown object")
var errObjectPropertyEffectObjectNotFound = errors.New("object property effect references an unknown object")
var errObjectPropertyEffectPropertyNotFound = errors.New("object property effect references an unknown property")
var errPropertyConditionPropertyNotFound = errors.New("property condition references an unknown property")
var errUseWithConditionObjectNotFound = errors.New("use with condition references an unknown object")
var errExitZoneNotFound = errors.New("zone exit references an unknown zone")
var errStartZoneNotFound = errors.New("start zone references an unknown zone")

func validateWorldState(ws WorldState) error {
	var errs []error
	if !ws.hasZone(ws.currentZoneId) {
		errs = append(errs, ReferenceError{
			kind: errStartZoneNotFound,
			Referent: ws.currentZoneId,
		})
	}
	errs = append(errs, validateExits(ws)...)
	errs = append(errs, validateObjects(ws)...)
	errs = append(errs, validateEffects(ws)...)
	errs = append(errs, validateConditionGroups(ws)...)
	return errors.Join(errs...)
}

func validateObjects(ws WorldState) []error {
	var errs []error
	for _, obj := range ws.objects {
		if !obj.inInventory() && !ws.hasZone(obj.zoneId) {
			errs = append(errs, ReferenceError{
				kind: errObjectZoneNotFound,
				Referrer: obj.id,
				Referent: obj.zoneId,
			})
		}
	}
	return errs
}

func validateEffects(ws WorldState) []error {
	var errs []error
	for _, o := range ws.objects {
		for _, a := range o.actions {
			for _, r := range a {
				for _, e := range r.effects {
					switch e := e.(type) {
					case CurrentZoneEffect:
						if !ws.hasZone(e.ZoneId) {
							errs = append(errs, ReferenceError{
								kind: errCurrentZoneEffectZoneNotFound,
								Referrer: o.id,
								Referent: e.ZoneId,
							})
						}
					case ObjectLocationEffect:
						if e.ZoneId != InventoryZoneId && !ws.hasZone(e.ZoneId) {
							errs = append(errs, ReferenceError{
								kind: errObjectLocationEffectZoneNotFound,
								Referrer: o.id,
								Referent: e.ZoneId,
							})
						}
						if !ws.hasObject(e.ObjectId) {
							errs = append(errs, ReferenceError{
								kind: errObjectLocationEffectObjectNotFound,
								Referrer: o.id,
								Referent: e.ObjectId,
							})
						}
					case ObjectTargetableEffect:
						if !ws.hasObject(e.ObjectId) {
							errs = append(errs, ReferenceError{
								kind: errObjectTargetableEffectObjectNotFound,
								Referrer: o.id,
								Referent: e.ObjectId,
							})
						}
					case ObjectVisibleEffect:
						if !ws.hasObject(e.ObjectId) {
							errs = append(errs, ReferenceError{
								kind: errObjectVisibleEffectObjectNotFound,
								Referrer: o.id,
								Referent: e.ObjectId,
							})
						}
					case ObjectPropertyEffect:
						obj, found := ws.object(e.ObjectId)
						if !found {
							errs = append(errs, ReferenceError{
								kind: errObjectPropertyEffectObjectNotFound,
								Referrer: o.id,
								Referent: e.ObjectId,
							})
						} else if !obj.hasProperty(e.PropertyId) {
							errs = append(errs, ReferenceError{
								kind: errObjectPropertyEffectPropertyNotFound,
								Referrer: o.id,
								Scope: e.ObjectId,
								Referent: e.PropertyId,
							})
						}
					}
				}
			}
		}
	}
	return errs
}

func validateConditionGroups(ws WorldState) []error {
	var errs []error
	for _, o := range ws.objects {
		for _, a := range o.actions {
			for _, r := range a {
				for _, c := range r.conditionGroups {
					if c.useWith != (Identity{}) && !ws.hasObject(c.useWith) {
						errs = append(errs, ReferenceError{
							kind: errUseWithConditionObjectNotFound,
							Referrer: o.id,
							Referent: c.useWith,
						})
					}
					for _, p := range c.properties {
						if !o.hasProperty(p.id) {
							errs = append(errs, ReferenceError{
								kind: errPropertyConditionPropertyNotFound,
								Referrer: o.id,
								Referent: p.id,
							})
						}
					}
				}
			}
		}
	}
	return errs
}

func validateExits(ws WorldState) []error {
	var errs []error
	for _, z := range ws.zones {
		for _, e := range z.exits {
			if !ws.hasZone(e.zoneId) {
				errs = append(errs, ReferenceError{
					kind: errExitZoneNotFound,
					Referrer: z.id,
					Referent: e.zoneId,
				})
			}
		}
	}
	return errs
}
