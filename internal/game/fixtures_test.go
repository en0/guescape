package game

import "testing"

type testOption[T any] func(*T)

func applyTestOptions[T any](v T, opts ...testOption[T]) T {
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

var testInventoryZoneId = InventoryZoneId.value

func testIdentity(t *testing.T, name string) Identity {
	t.Helper()
	id, err := NewIdentity(name)
	if err != nil {
		t.Fatalf("NewIdentity(%q): %v", name, err)
	}
	return id
}

func testIdentities(t *testing.T, names ...string) []Identity {
	t.Helper()
	var result []Identity
	for _, name := range names {
		id, err := NewIdentity(name)
		if err != nil {
			t.Fatalf("NewIdentity(%q): %v", name, err)
		}
		result = append(result, id)
	}
	return result
}

// WorldState

func testWorldState(
	t *testing.T,
	startZoneId string,
	desc string,
	zones []Zone,
	objects []Object,
	opts ...testOption[WorldState],
) WorldState {
	t.Helper()
	szid := testIdentity(t, startZoneId)
	ws, err := NewWorldState(szid, desc, nil, zones, objects, nil)
	if err != nil {t.Fatalf("NewWorldState(...): %v", err)}
	return applyTestOptions(ws, opts...)
}

func testWithWorldStateFlag(t *testing.T, flagId string) testOption[WorldState] {
	t.Helper()
	id := testIdentity(t, flagId)
	return func(ws *WorldState) {
		ws.flags = append(ws.flags, id)
	}
}

// Object

func testObject(t *testing.T, id string, zoneId string, opts ...testOption[Object]) Object {
	t.Helper()
	return applyTestOptions(NewObject(
		testIdentity(t, id),
		testIdentity(t, zoneId),
		id, nil, nil, nil, nil, nil,
	), opts...)
}

func testWithObjectDesc(d string) testOption[Object] {
	return func(o *Object) {o.desc = d}
}

func testWithObjectVisibility(v bool) testOption[Object] {
	return func(o *Object) {o.visible = v}
}

func testWithObjectTargetability(v bool) testOption[Object] {
	return func(o *Object) {o.targetable = v}
}

func testWithObjectCarriability(v bool) testOption[Object] {
	return func(o *Object) {o.carriable = v}
}

func testWithObjectProperty(t *testing.T, prop string, value string) testOption[Object] {
	t.Helper()
	n := testIdentity(t, prop)
	return func(o *Object) {
		if o.properties == nil {
			o.properties = make(map[Identity]string)
		}
		o.properties[n] = value
	}
}

func testWithObjectAction(v ObjectVerb, r ...Rule) testOption[Object] {
	return func(o *Object) {
		if o.actions == nil {
			o.actions = make(map[ObjectVerb][]Rule)
		}
		o.actions[v] = r
	}
}

// Rule

func testRule(t *testing.T, effects []Effect, opts ...testOption[Rule]) Rule {
	t.Helper()
	return applyTestOptions(NewRule(nil, effects), opts...)
}

func testWithConditionGroups(cg ...ConditionGroup) testOption[Rule] {
	return func(r *Rule) {r.conditionGroups = append(r.conditionGroups, cg...)}
}

// Zone

func testZone(t *testing.T, id string, opts ...testOption[Zone]) Zone {
	t.Helper()
	return applyTestOptions(NewZone(testIdentity(t, id), id, nil), opts...)
}

func testWithZoneExit(t *testing.T, key string, zoneId string, flags ...string) testOption[Zone] {
	t.Helper()
	kid := testIdentity(t, key)
	zid := testIdentity(t, zoneId)
	flgs := testIdentities(t, flags...)
	return func(z *Zone) {
		if z.exits == nil {z.exits = make(map[Identity]ZoneExit)}
		z.exits[kid] = NewZoneExit(zid, flgs)
	}
}

// ConditionGroup

func testConditionGroup(t *testing.T, opts ...testOption[ConditionGroup]) ConditionGroup {
	t.Helper()
	return applyTestOptions(NewConditionGroup(Identity{}, nil, nil), opts...)
}

func testWithUseWithCondition(t *testing.T, item string) testOption[ConditionGroup] {
	t.Helper()
	id := testIdentity(t, item)
	return func(cg *ConditionGroup) {cg.useWith = id}
}

func testWithPropertyCondition(t *testing.T, prop string, value string) testOption[ConditionGroup] {
	t.Helper()
	id := testIdentity(t, prop)
	return func(cg *ConditionGroup) {
		cg.properties = append(cg.properties, NewPropertyCondition(id, value))
	}
}

func testWithFlagCondition(t *testing.T, flag string, value bool) testOption[ConditionGroup] {
	t.Helper()
	id := testIdentity(t, flag)
	return func(cg *ConditionGroup) {
		cg.flags = append(cg.flags, NewFlagCondition(id, value))
	}
}
