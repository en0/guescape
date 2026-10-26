package game

import (
	"errors"
	"testing"
)

type validateCase struct {
	name string
	err error
	state WorldState
}

func runValidateCases(t *testing.T, fn func(WorldState) []error, cases []validateCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Join(fn(tt.state)...)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Unexpected result: err=%v, wants err=%v", err, tt.err)
			}
		})
	}
}


func TestValidateWorldState(t *testing.T) {
	w := func (ws WorldState) []error {return []error{validateWorldState(ws)}}
	runValidateCases(t, w, []validateCase{
		{
			name: "valid-state",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				desc: "hello-world",
				flags: testIdentities(t, "flag-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "east", "zone-two", "flag-one")),
					testZone(t, "zone-two"),
				},
				objects: []Object{
					testObject(t, "obj-one", "zone-one",
						testWithObjectProperty(t, "prop-one", "prop-one-val"),
					),
					testObject(t, "obj-two", "zone-two",
						testWithObjectAction(VerbUse,
							testRule(t,
								[]Effect{
									SayEffect{"set prop on obj-one"},
									ObjectPropertyEffect{
										ObjectId: testIdentity(t, "obj-one"),
										PropertyId: testIdentity(t, "prop-one"),
										Value: "prop-one-val2",
									},
								},
								testWithConditionGroups(
									testConditionGroup(t,
										testWithFlagCondition(t, "flag-one", true),
									),
								),
							),
						),
					),
				},
			},
		},
		{
			name: "start-zone-no-exist",
			err: errStartZoneNotFound,
			state: WorldState{
				currentZoneId: testIdentity(t, "no-exists"),
				zones: []Zone{},
			},
		},
	})
}

func TestValidateObjects(t *testing.T) {
	runValidateCases(t, validateObjects, []validateCase{
		{
			name: "zone-exists",
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{testObject(t, "obj-one", "zone-one")},
			},
		},
		{
			name: "zone-no-exist",
			err: errObjectZoneNotFound,
			state: WorldState{
				objects: []Object{testObject(t, "obj-one", "zone-no-exist")},
			},
		},
		{
			name: "zone--can-be-inventory",
			state: WorldState{
				objects: []Object{testObject(t, "obj-one", testInventoryZoneId)},
			},
		},
	})
}

func TestValidateEffects(t *testing.T) {
	runValidateCases(t, validateEffects, []validateCase{

		// CurrentZoneEffect Tests
		{
			name: "current-zone-effect-zone-no-exist",
			err: errCurrentZoneEffectZoneNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								CurrentZoneEffect{testIdentity(t, "zone-no-exist")},
							}),
						),
					),
				},
			},
		},

		// ObjectLocationEffect Tests
		{
			name: "object-location-effect-zone-no-exist",
			err: errObjectLocationEffectZoneNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectLocationEffect{
									ObjectId: testIdentity(t, "obj-one"),
									ZoneId: testIdentity(t, "zone-no-exist"),
								},
							}),
						),
					),
				},
			},
		},
		{
			name: "object-location-effect-can-be-inventory",
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectLocationEffect{
									ObjectId: testIdentity(t, "obj-one"),
									ZoneId: testIdentity(t, testInventoryZoneId),
								},
							}),
						),
					),
				},
			},
		},
		{
			name: "object-location-effect-object-must-exist",
			err: errObjectLocationEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectLocationEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
									ZoneId: testIdentity(t, "zone-one"),
								},
							}),
						),
					),
				},
			},
		},

		// ObjectTargetableEffect Tests
		{
			name: "object-targetable-effect-object-must-exist",
			err: errObjectTargetableEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectTargetableEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
								},
							}),
						),
					),
				},
			},
		},

		// ObjectVisibleEffect Tests
		{
			name: "object-visible-effect-object-must-exist",
			err: errObjectVisibleEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectVisibleEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
								},
							}),
						),
					),
				},
			},
		},

		// ObjectPropertyEffect Tests
		{
			name: "object-property-effect-object-must-exist",
			err: errObjectPropertyEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectProperty(t, "prop-one", "prop-one-val"),
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectPropertyEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
									PropertyId: testIdentity(t, "prop-one"),
								},
							}),
						),
					),
				},
			},
		},
		{
			name: "object-property-effect-property-must-exist",
			err: errObjectPropertyEffectPropertyNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectProperty(t, "prop-one", "prop-one-val"),
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectPropertyEffect{
									ObjectId: testIdentity(t, "obj-one"),
									PropertyId: testIdentity(t, "prop-no-exist"),
								},
							}),
						),
					),
				},
			},
		},
	})
}

func TestValidateConditionGroups(t *testing.T) {
	runValidateCases(t, validateConditionGroups, []validateCase{

		// PropertyCondition Tests
		{
			name: "property-condition-property-no-exist",
			err: errPropertyConditionPropertyNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{SayEffect{""}},
								testWithConditionGroups(
									testConditionGroup(t,
										testWithPropertyCondition(t, "prop-no-exist", ""),
									),
								),
							),
						),
					),
				},
			},
		},

		// useWith Tests
		{
			name: "use-with-condition-object-no-exist",
			err: errUseWithConditionObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{SayEffect{""}},
								testWithConditionGroups(
									testConditionGroup(t,
										testWithUseWithCondition(t, "item-no-exist"),
									),
								),
							),
						),
					),
				},
			},
		},
		{
			name: "use-with-condition-object-can-be-empty",
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{SayEffect{""}},
								testWithConditionGroups(
									testConditionGroup(t),
								),
							),
						),
					),
				},
			},
		},
	})
}

func TestValidateExits(t *testing.T) {
	runValidateCases(t, validateExits, []validateCase{
		{
			name: "exit-zone-no-exist",
			err: errExitZoneNotFound,
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "east", "zone-no-exist")),
				},
			},
		},
	})
}
