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
						testWithObjectState(t, "state-one", "state-one-val"),
					),
					testObject(t, "obj-two", "zone-two",
						testWithObjectAction(VerbUse,
							testRule(t,
								[]Effect{
									SayEffect{"set state on obj-one"},
									SetObjectStateEffect{
										ObjectId: testIdentity(t, "obj-one"),
										StateId: testIdentity(t, "state-one"),
										Value: "state-one-val2",
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

		// ObjectCarriableEffect Tests
		{
			name: "object-carriable-effect-object-must-exist",
			err: errObjectCarriableEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ObjectCarriableEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
								},
							}),
						),
					),
				},
			},
		},

		// ObjectStateEffect Tests
		{
			name: "object-state-effect-object-must-exist",
			err: errObjectStateEffectObjectNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectState(t, "state-one", "state-one-val"),
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								SetObjectStateEffect{
									ObjectId: testIdentity(t, "obj-no-exist"),
									StateId: testIdentity(t, "state-one"),
								},
							}),
						),
					),
				},
			},
		},
		{
			name: "object-state-effect-state-must-exist",
			err: errObjectStateEffectStateNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectState(t, "state-one", "state-one-val"),
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								SetObjectStateEffect{
									ObjectId: testIdentity(t, "obj-one"),
									StateId: testIdentity(t, "state-no-exist"),
								},
							}),
						),
					),
				},
			},
		},

		// ActivateModeEffect

		{
			name: "activate-mode-effect-mode-must-exist",
			err: errModeEffectModeNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								ActivateModeEffect{
									ModeId: testIdentity(t, "mode-no-exist"),
								},
							}),
						),
					),
				},
			},
		},

		// SetModeStateEffect

		{
			name: "set-mode-state-effect-mode-must-exist",
			err: errModeStateEffectModeNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								SetModeStateEffect{
									ModeId: testIdentity(t, "mode-no-exist"),
									State: testIdentity(t, "state-one"),
									Value: "state-one-value",
								},
							}),
						),
					),
				},
			},
		},
		{
			name: "set-mode-state-effect-state-must-exist",
			err: errModeStateEffectStateNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				modes: []Mode{testMode(t, "mode-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								SetModeStateEffect{
									ModeId: testIdentity(t, "mode-one"),
									State: testIdentity(t, "state-one"),
									Value: "state-one-value",
								},
							}),
						),
					),
				},
			},
		},

		// SetModePromptEffect

		{
			name: "set-mode-prompt-effect-mode-must-exist",
			err: errModePromptEffectModeNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{
								SetModePromptEffect{
									ModeId: testIdentity(t, "mode-no-exist"),
									Value: "prompt-value",
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

		// ObjectStateCondition Tests
		{
			name: "object-state-condition-state-no-exist",
			err: errObjectStateConditionStateNotFound,
			state: WorldState{
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "obj-one", "zone-one",
						testWithObjectAction(VerbUse,
							testRule(t, []Effect{SayEffect{""}},
								testWithConditionGroups(
									testConditionGroup(t,
										testWithObjectStateCondition(t, "state-no-exist", ""),
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
