package game

import (
	"errors"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {

	newWorldState := func () WorldState {
		return testWorldState(t, "zone-one", "desc",
			[]Zone{
				testZone(t, "zone-exists"),
				testZone(t, "zone-one", testWithZoneExit(t, "exit-one", "zone-exists")),
			},
			[]Object{
				testObject(t, "obj-exists", "zone-one",
					testWithObjectVisibility(true),
					testWithObjectTargetability(true),
					testWithObjectCarriability(true),
					testWithObjectProperty(t, "prop-exists", "old-val"),
				),
			},
			testWithWorldStateModes(t, true,
				testMode(t, "mode-one",
					testWithModeState(t, "state-one", "initial"),
					testWithModePrompt("initial"),
				),
				testMode(t, "mode-two"),
			),
			testWithWorldStateFlag(t, "flag-exists"),
		)
	}

	tests := []struct {
		name string
		err error
		ef Effect
		check func(WorldState, error) bool
	} {
		{
			name: "ignores-say-effect",
			ef: SayEffect{Text: "hello"},
			check: func(ws WorldState, err error) bool {
				return reflect.DeepEqual(ws, newWorldState())
			},
		},
		{
			name: "ignores-end-game-effect",
			ef: EndGameEffect{IsWin: false, Text: "hello"},
			check: func(ws WorldState, err error) bool {
				return reflect.DeepEqual(ws, newWorldState())
			},
		},
		{
			name: "current-zone-effect-changes-current-zone",
			ef: CurrentZoneEffect{ZoneId: testIdentity(t, "zone-exists")},
			check: func(ws WorldState, err error) bool {
				return ws.currentZoneId.value == "zone-exists"
			},
		},
		{
			name: "add-flag-effect-adds-missing-flag",
			ef: AddFlagEffect{Flag: testIdentity(t, "flag-no-exists")},
			check: func(ws WorldState, err error) bool {
				return ws.HasFlag(testIdentity(t, "flag-no-exists"))
			},
		},
		{
			name: "add-flag-effect-skips-if-flag-exists",
			ef: AddFlagEffect{Flag: testIdentity(t, "flag-exists")},
			check: func(ws WorldState, err error) bool {
				return len(filter(ws.flags, func(f Identity) bool { return f.value == "flag-exists"
				})) == 1
			},
		},
		{
			name: "remove-flag-effect-removes-existing-flag",
			ef: RemoveFlagEffect{Flag: testIdentity(t, "flag-exists")},
			check: func(ws WorldState, err error) bool {
				return !ws.HasFlag(testIdentity(t, "flag-exists"))
			},
		},
		{
			name: "remove-flag-effect-skips-missing-flag",
			ef: RemoveFlagEffect{Flag: testIdentity(t, "flag-no-exists")},
			check: func(ws WorldState, err error) bool {
				return len(ws.flags) == 1
			},
		},
		{
			name: "object-property-effect-sets-property-value",
			ef: ObjectPropertyEffect{
				ObjectId: testIdentity(t, "obj-exists"),
				PropertyId: testIdentity(t, "prop-exists"),
				Value: "new-value",
			},
			check: func(ws WorldState, err error) bool {
				obj, found := ws.object(testIdentity(t, "obj-exists"))
				if !found {return false}
				prop, found := obj.property(testIdentity(t, "prop-exists"))
				if !found {return false}
				return prop == "new-value"
			},
		},
		{
			name: "object-property-effect-targeting-missing-object-is-an-error",
			err: errEffectObjectNotFound,
			ef: ObjectPropertyEffect{
				ObjectId: testIdentity(t, "obj-no-exists"),
				PropertyId: testIdentity(t, "prop-exists"),
				Value: "new-value",
			},
		},
		{
			name: "object-property-effect-targeting-missing-property-is-an-error",
			err: errEffectPropertyNotFound,
			ef: ObjectPropertyEffect{
				ObjectId: testIdentity(t, "obj-exists"),
				PropertyId: testIdentity(t, "prop-no-exists"),
				Value: "new-value",
			},
		},
		{
			name: "object-visible-effect-sets-object-visibility",
			ef: ObjectVisibleEffect{ObjectId: testIdentity(t, "obj-exists"), Visible: false},
			check: func(ws WorldState, err error) bool {
				o, found := ws.object(testIdentity(t, "obj-exists"))
				if !found { return false }
				return o.visible == false
			},
		},
		{
			name: "object-visible-effect-targeting-missing-object-is-an-error",
			err: errEffectObjectNotFound,
			ef: ObjectVisibleEffect{ObjectId: testIdentity(t, "obj-no-exists"), Visible: false},
		},
		{
			name: "object-targetable-effect-sets-object-targetability",
			ef: ObjectTargetableEffect{ObjectId: testIdentity(t, "obj-exists"), Targetable: false},
			check: func(ws WorldState, err error) bool {
				o, found := ws.object(testIdentity(t, "obj-exists"))
				if !found { return false }
				return o.targetable == false
			},
		},
		{
			name: "object-targetable-effect-targeting-missing-object-is-an-error",
			err: errEffectObjectNotFound,
			ef: ObjectTargetableEffect{ObjectId: testIdentity(t, "obj-no-exists"), Targetable: false},
		},
		{
			name: "object-carriable-effect-sets-object-carriability",
			ef: ObjectCarriableEffect{ObjectId: testIdentity(t, "obj-exists"), Carriable: false},
			check: func(ws WorldState, err error) bool {
				o, found := ws.object(testIdentity(t, "obj-exists"))
				if !found { return false }
				return o.carriable == false
			},
		},
		// {
		// 	name: "object-targetable-effect-targeting-missing-object-is-an-error",
		// 	err: errEffectObjectNotFound,
		// 	ef: ObjectTargetableEffect{ObjectId: testIdentity(t, "obj-no-exists"), Targetable: false},
		// },
		{
			name: "object-location-effect-sets-object-location",
			ef: ObjectLocationEffect{
				ObjectId: testIdentity(t, "obj-exists"),
				ZoneId: testIdentity(t, "zone-exists"),
			},
			check: func(ws WorldState, err error) bool {
				o, found := ws.object(testIdentity(t, "obj-exists"))
				if !found { return false }
				return o.zoneId.value == "zone-exists"
			},
		},
		{
			name: "object-location-effect-targeting-missing-object-is-an-errors",
			err: errEffectObjectNotFound,
			ef: ObjectLocationEffect{
				ObjectId: testIdentity(t, "obj-no-exists"),
				ZoneId: testIdentity(t, "zone-exists"),
			},
		},
		{
			name: "set-mode-effect-sets-current-mode-id",
			ef: SetModeEffect{
				ModeId: testIdentity(t, "mode-two"),
			},
			check: func(ws WorldState, err error) bool {
				return ws.currentModeId == testIdentity(t, "mode-two")
			},
		},
		{
			name: "clear-mode-effect-sets-current-mode-id-to-zero-identity",
			ef: ClearModeEffect{},
			check: func(ws WorldState, err error) bool {
				return ws.currentModeId == (Identity{})
			},
		},
		{
			name: "set-mode-state-effect-sets-the-mode-state-when-present",
			ef: SetModeStateEffect{
				ModeId: testIdentity(t, "mode-one"),
				State: testIdentity(t, "state-one"),
				Value: "foo",
			},
			check: func(ws WorldState, err error) bool {
				m, found := ws.mode(testIdentity(t, "mode-one"))
				if !found { return false }
				if v, ok := m.states[testIdentity(t, "state-one")]; ok {
					return v == "foo"
				}
				return false
			},
		},
		{
			name: "set-mode-state-effect-returns-error-when-the-state-is-not-present",
			ef: SetModeStateEffect{
				ModeId: testIdentity(t, "mode-one"),
				State: testIdentity(t, "state-no-exists"),
				Value: "foo",
			},
			err: errEffectStateNotFound,
		},
		{
			name: "set-mode-state-effect-returns-error-when-the-mode-is-not-present",
			ef: SetModeStateEffect{
				ModeId: testIdentity(t, "mode-no-exists"),
				State: testIdentity(t, "state-one"),
				Value: "foo",
			},
			err: errEffectModeNotFound,
		},
		{
			name: "set-mode-prompt-effect-sets-mode-prompt",
			ef: SetModePromptEffect{
				ModeId: testIdentity(t, "mode-one"),
				Value: "foo",
			},
			check: func(ws WorldState, err error) bool {
				m, found := ws.mode(testIdentity(t, "mode-one"))
				if !found { return false }
				return m.prompt == "foo"
			},
		},
		{
			name: "set-mode-prompt-effect-returns-error-when-mode-not-present",
			ef: SetModePromptEffect{
				ModeId: testIdentity(t, "mode-no-exists"),
				Value: "foo",
			},
			err: errEffectModeNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsIn := newWorldState()
			wsOut, err := Apply(wsIn, tt.ef)
			if !reflect.DeepEqual(wsIn, newWorldState()) {
				t.Fatal("Apply Mutation detected!")
			} else if !errors.Is(err, tt.err) {
				t.Fatalf("Apply: unexpected error: %v", err)
			} else if tt.check != nil && !tt.check(wsOut, err) {
				t.Fatalf("Apply: check failed: err=%v, ws=%v", err, wsOut)
			} else if tt.check == nil && tt.err == nil {
				t.Fatalf("test case asserts nothing")
			}
		})
	}
}

