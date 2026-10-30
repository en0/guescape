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
