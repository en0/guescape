package game

import (
	"reflect"
	"testing"
)

// TODO: rebuild tests, focus on default values and presents of factory-provided components.
func TestNewWorldState(t *testing.T) {
	description := "test-new-world-state"
	tests := []struct {
		name string
		zone string
		zones []Zone
		objs []Object
		check func(WorldState, error) bool
	} {
		{
			name: "create-valid-world",
			zone: "zone-one",
			zones: []Zone{testZone(t, "zone-one")},
			check: func(ws WorldState, err error) bool {
				return err == nil
			},
		},
		{
			// This is an indirect test. Just want to make sure
			// validate is called in world construction
			name: "create-world-calls-validate",
			zone: "zone-no-exist",
			zones: []Zone{testZone(t, "zone-one")},
			check: func(ws WorldState, err error) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startZoneId := testIdentity(t, tt.zone)
			state, err := NewWorldState(
				startZoneId, description,
				nil, tt.zones, tt.objs, nil, nil,
			)
			if !tt.check(state, err) {
				t.Fatalf("Unexpected result: err=%v", err)
			}
		})
	}
}

func TestInventory(t *testing.T) {

	wolrdStateFactory := func(o ...Object) WorldState {
		return testWorldState(t,
			"zone-one", "world-desc",
			[]Zone{testZone(t, "zone-one")},
			o,
		)
	}

	tests := []struct{
		name string
		ws WorldState
		expect []Identity
	} {
		{
			name: "inventory-with-objects",
			ws: wolrdStateFactory(
				testObject(t, "obj-one", testInventoryZoneId),
				testObject(t, "obj-two", testInventoryZoneId),
				testObject(t, "obj-three", "zone-one"),
			),
			expect: testIdentities(t, "obj-one", "obj-two"),
		},
		{
			name: "inventory-empty",
			ws: wolrdStateFactory(),
			expect: nil,
		},
		{
			name: "inventory-empty-but-world-objects",
			ws: wolrdStateFactory(
				testObject(t, "obj-two", "zone-one"),
				testObject(t, "obj-one", "zone-one"),
			),
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.ws.Inventory()
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Fatalf("ws.Inventory() %v, want %v", actual, tt.expect)
			}
		})
	}
}

func TestViewCurrentZone(t *testing.T) {
	tests := []struct{
		name string
		ws WorldState
		expect ZoneView
	} {
		{
			name: "empty-world",
			ws: testWorldState(t, "zone-one", "world-desc",
				[]Zone{testZone(t, "zone-one")},
				nil,
			),
			expect: ZoneView{
				Id: testIdentity(t, "zone-one"),
				Desc: "zone-one",
				Exits: nil,
				Objects: nil,
			},
		},
		{
			name: "world-zone-exits",
			ws: testWorldState(t, "zone-one", "world-desc",
				[]Zone{
					testZone(t, "zone-two"),
					testZone(t, "zone-one",
						testWithZoneExit(t, "east", "zone-two"),
					),
				},
				[]Object{
				},
			),
			expect: ZoneView{
				Id: testIdentity(t, "zone-one"),
				Desc: "zone-one",
				Exits: testIdentities(t, "east"),
				Objects: nil,
			},
		},
		{
			name: "empty-world",
			ws: testWorldState(t, "zone-one", "world-desc",
				[]Zone{testZone(t, "zone-one")},
				[]Object{
					testObject(t, "obj-one", "zone-one", testWithObjectVisibility(true)),
					testObject(t, "obj-two", "zone-one", testWithObjectVisibility(true)),
				},
			),
			expect: ZoneView{
				Id: testIdentity(t, "zone-one"),
				Desc: "zone-one",
				Exits: nil,
				Objects: testIdentities(t, "obj-one", "obj-two"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.ws.ViewCurrentZone()
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Fatalf("ws.ViewCurrentZone() %v, want %v", actual, tt.expect)
			}
		})
	}
}

func TestPrompt(t *testing.T) {
	tests := []struct{
		name string
		ws WorldState
		expect PromptView
	}{
		{
			name: "prompt-view-contains-current-zone",
			ws: testWorldState(t, "zone-one", "zone-one", []Zone{testZone(t, "zone-one")}, nil),
			expect: PromptView{Text: "zone-one> ", ModeActive: false},
		},
		{
			name: "prompt-view-shows-active-mode",
			ws: testWorldState(t, "zone-one", "zone-one",
				[]Zone{testZone(t, "zone-one")},
				nil,
				testWithWorldStateModes(t, true, testMode(t, "mode-one")),
			),
			expect: PromptView{Text: "", ModeActive: true},
		},
		{
			name: "prompt-view-contains-mode-prompt",
			ws: testWorldState(t, "zone-one", "zone-one",
				[]Zone{testZone(t, "zone-one")},
				nil,
				testWithWorldStateModes(t, true, testMode(t, "mode-one", testWithModePrompt("test-one"))),
			),
			expect: PromptView{Text: "test-one", ModeActive: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.ws.Prompt()
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Fatalf("ws.Prompt() %v, want %v", actual, tt.expect)
			}
		})
	}
}

func TestHasFlag(t *testing.T) {
	tests := []struct{
		name string
		flag string
		expect bool
		ws WorldState
	}{
		{
			name: "has-flag-returns-true-when-flag-present",
			flag: "flag-one",
			expect: true,
			ws: WorldState{
				flags: testIdentities(t, "flag-one"),
			},
		},
		{
			name: "has-flag-returns-false-when-flag-not-present",
			flag: "flag-one",
			expect: false,
			ws: WorldState{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := testIdentity(t, tt.flag)
			actual := tt.ws.HasFlag(id)
			if actual != tt.expect {
				t.Fatalf("ws.HasFlag(%v) %v, want %v", tt.flag, actual, tt.expect)
			}
		})
	}
}

func TestObjectProperty(t *testing.T) {
	tests := []struct{
		name string
		obj string
		property string
		value string
		found bool
		ws WorldState
	}{
		{
			name: "object-property-returns-property-value-when-present",
			obj: "obj-one",
			property: "prop-one",
			value: "prop-one-value",
			found: true,
			ws: WorldState{
				objects: []Object{
					testObject(t, "obj-one", "zone-one", testWithObjectProperty(t, "prop-one", "prop-one-value")),
				},
			},
		},
		{
			name: "object-property-returns-not-oufnd-when-property-not-present",
			obj: "obj-one",
			property: "prop-one",
			value: "",
			found: false,
			ws: WorldState{
				objects: []Object{
					testObject(t, "obj-one", "zone-one"),
				},
			},
		},
		{
			name: "object-property-returns-not-found-when-property-not-present",
			obj: "obj-one",
			property: "prop-one",
			value: "",
			found: false,
			ws: WorldState{
				objects: []Object{
					testObject(t, "obj-one", "zone-one"),
				},
			},
		},
		{
			name: "object-property-returns-not-found-when-object-not-present",
			obj: "obj-one",
			property: "prop-one",
			value: "",
			found: false,
			ws: WorldState{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objId := testIdentity(t, tt.obj)
			propId := testIdentity(t, tt.property)
			actual, found := tt.ws.ObjectProperty(objId, propId)
			if found != tt.found || actual != tt.value {
				t.Fatalf("ws.ObjectProperty(%v, %v) (%v, %v) want (%v, %v)", tt.obj, tt.property, actual, found, tt.value, tt.found)
			}
		})
	}
}

func TestCurrentModeState(t *testing.T) {
	tests := []struct{
		name string
		state string
		value string
		found bool
		ws WorldState
	}{
		{
			name: "current-mode-state-returns-state-value-when-present",
			state: "state-one",
			value: "state-one-value",
			found: true,
			ws: WorldState{
				currentModeId: testIdentity(t, "mode-one"),
				modes: []Mode{
					testMode(t, "mode-one", testWithModeState(t, "state-one", "state-one-value")),
				},
			},
		},
		{
			name: "current-mode-state-returns-not-found-when-state-not-present",
			state: "state-one",
			found: false,
			ws: WorldState{
				currentModeId: testIdentity(t, "mode-one"),
				modes: []Mode{
					testMode(t, "mode-one"),
				},
			},
		},
		{
			name: "current-mode-state-returns-not-found-when-mode-not-present",
			state: "state-one",
			found: false,
			ws: WorldState{
				currentModeId: testIdentity(t, "mode-one"),
			},
		},
		{
			name: "current-mode-state-returns-not-found-when-mode-not-set",
			state: "state-one",
			found: false,
			ws: WorldState{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := testIdentity(t, tt.state)
			actual, found := tt.ws.CurrentModeState(id)
			if found != tt.found || actual != tt.value {
				t.Fatalf("ws.CurrentModeState(%v) (%v, %v) want (%v, %v)", tt.state, actual, found, tt.value, tt.found)
			}
		})
	}
}
