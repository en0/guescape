package game

import (
	"slices"
	"testing"
)

type effectCase struct {
	name string
	state WorldState
	action Action
	expected []Effect
}

func runEffectCases(t *testing.T, fn func(WorldState, Action) []Effect, tests []effectCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := fn(tt.state, tt.action)
			if !slices.Equal(tt.expected, actual) {
				t.Errorf("effect types = %v want %v", actual, tt.expected)
			}
		})
	}
}

func TestPlanEffects(t *testing.T) {
	runEffectCases(t, PlanEffects, []effectCase{
		{
			name: "move-zone-unblocked",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "zone-exit", "zone-two")),
					testZone(t, "zone-two"),
				},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{CurrentZoneEffect{testIdentity(t, "zone-two")}},
		}, {
			name: "move-zone-has-flags",
			state: WorldState{
				flags: testIdentities(t, "exit-flag"),
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "zone-exit", "zone-two", "exit-flag")),
					testZone(t, "zone-two"),
				},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{CurrentZoneEffect{testIdentity(t, "zone-two")}},
		}, {
			name: "move-zone-missing-flags",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "zone-exit", "zone-two", "exit-flag")),
					testZone(t, "zone-two"),
				},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{SayEffect{msgDefaults[MsgWayBlocked]}},
		}, {
			name: "move-zone-missing-flags-custom-message",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				messages: map[MsgKey]string{MsgWayBlocked: "way-blocked"},
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "zone-exit", "zone-two", "exit-flag")),
					testZone(t, "zone-two"),
				},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{SayEffect{"way-blocked"}},
		}, {
			name: "move-no-exit",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one"), testZone(t, "zone-two")},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{SayEffect{msgDefaults[MsgCannotMove]}},
		}, {
			name: "move-bad-zone",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{
					testZone(t, "zone-one", testWithZoneExit(t, "zone-exit", "zone-two")),
				},
			},
			action: NewNavAction(CmdMove, testIdentity(t, "zone-exit")),
			expected: []Effect{SayEffect{msgDefaults[MsgDestinationMissing]}},
		}, {
			name: "examine-object",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{testObject(t, "object-one", "zone-one")},
			},
			action: NewObjectAction(VerbExamine, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"object-one"}},
		}, {
			name: "examine-object-not-found",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
			},
			action: NewObjectAction(VerbExamine, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotFound]}},
		}, {
			name: "examine-object-not-found-custom-message",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{},
				messages: map[MsgKey]string{MsgObjectNotFound: "not-found"},
			},
			action: NewObjectAction(VerbExamine, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"not-found"}},
		}, {
			name: "examine-object-in-inventory",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", testInventoryZoneId),
				},
				messages: map[MsgKey]string{MsgObjectNotFound: "not-found"},
			},
			action: NewObjectAction(VerbExamine, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"object-one"}},
		}, {
			name: "take-object",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectCarriability(true)),
				},
			},
			action: NewObjectAction(VerbTake, testIdentity(t, "object-one")),
			expected: []Effect{ObjectLocationEffect{
				ObjectId: testIdentity(t, "object-one"),
				ZoneId: InventoryZoneId,
			}},
		}, {
			name: "take-object-not-found",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
			},
			action: NewObjectAction(VerbTake, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotFound]}},
		}, {
			name: "take-object-already-taken",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", testInventoryZoneId, testWithObjectCarriability(true)),
				},
			},
			action: NewObjectAction(VerbTake, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectAlreadyTaken]}},
		}, {
			name: "take-object-not-carriable",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectCarriability(false)),
				},
			},
			action: NewObjectAction(VerbTake, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotCarriable]}},
		}, {
			name: "drop-object",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", testInventoryZoneId),
				},
			},
			action: NewObjectAction(VerbDrop, testIdentity(t, "object-one")),
			expected: []Effect{ObjectLocationEffect{
				ObjectId: testIdentity(t, "object-one"),
				ZoneId: testIdentity(t, "zone-one"),
			}},
		}, {
			name: "drop-object-not-found",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
			},
			action: NewObjectAction(VerbDrop, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotFound]}},
		}, {
			name: "drop-object-not-in-inventory",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one"),
				},
			},
			action: NewObjectAction(VerbDrop, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotTaken]}},
		}, {
			name: "use-object",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one"),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotUsable]}},
		}, {
			name: "use-object-not-found",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgObjectNotFound]}},
		}, {
			name: "use-object-with-item-not-found",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one"),
				},
			},
			action: NewObjectActionWithItem(
				VerbUse,
				testIdentity(t, "object-one"),
				testIdentity(t, "object-two"),
			),
			expected: []Effect{SayEffect{msgDefaults[MsgItemNotFound]}},
		}, {
			name: "use-object-with-item-no-handler",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one"),
					testObject(t, "object-two", testInventoryZoneId),
				},
			},
			action: NewObjectActionWithItem(
				VerbUse,
				testIdentity(t, "object-one"),
				testIdentity(t, "object-two"),
			),
			expected: []Effect{SayEffect{msgDefaults[MsgNoHandler]}},
		}, {
			name: "custom-use",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t, []Effect{SayEffect{"custom-action"}}),
					)),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"custom-action"}},
		}, {
			name: "custom-use-with-flag-condition-unmet",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t,
							[]Effect{SayEffect{"custom-action"}},
							testWithConditionGroups(testConditionGroup(t,
								testWithFlagCondition(t, "flag-one", true),
							)),
						),
					)),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgNoHandler]}},
		}, {
			name: "custom-use-with-flag-condition-met",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				flags: []Identity{testIdentity(t, "flag-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t,
							[]Effect{SayEffect{"flag-present"}},
							testWithConditionGroups(testConditionGroup(t,
								testWithFlagCondition(t, "flag-one", true),
							)),
						),
					)),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"flag-present"}},
		}, {
			name: "custom-use-with-property-condition-unmet",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t,
							[]Effect{SayEffect{"custom-action"}},
							testWithConditionGroups(testConditionGroup(t,
								testWithPropertyCondition(t, "prop-one", "foo"),
							)),
						),
					)),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgNoHandler]}},
		}, {
			name: "custom-use-with-property-condition-met",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(
						t, "object-one", "zone-one",
						testWithObjectProperty(t, "prop-one", "foo"),
						testWithObjectAction(VerbUse,
							testRule(t,
								[]Effect{SayEffect{"has-property-foo"}},
								testWithConditionGroups(testConditionGroup(t,
									testWithPropertyCondition(t, "prop-one", "foo"),
								)),
							),
						),
					),
				},
			},
			action: NewObjectAction(VerbUse, testIdentity(t, "object-one")),
			expected: []Effect{SayEffect{"has-property-foo"}},
		}, {
			name: "custom-use-with-item-condition-unmet",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "item-one", testInventoryZoneId),
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t,
							[]Effect{SayEffect{"custom-action"}},
							testWithConditionGroups(testConditionGroup(t,
								testWithUseWithCondition(t, "no-item"),
							)),
						),
					)),
				},
			},
			action: NewObjectActionWithItem(VerbUse, testIdentity(t, "object-one"), testIdentity(t, "item-one")),
			expected: []Effect{SayEffect{msgDefaults[MsgNoHandler]}},
		}, {
			name: "custom-use-with-item-condition-met",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				objects: []Object{
					testObject(t, "item-one", testInventoryZoneId),
					testObject(t, "object-one", "zone-one", testWithObjectAction(VerbUse,
						testRule(t,
							[]Effect{SayEffect{"has-item-one"}},
							testWithConditionGroups(testConditionGroup(t,
								testWithUseWithCondition(t, "item-one"),
							)),
						),
					)),
				},
			},
			action: NewObjectActionWithItem(VerbUse, testIdentity(t, "object-one"), testIdentity(t, "item-one")),
			expected: []Effect{SayEffect{"has-item-one"}},
		}, {
			name: "custom-use-with-multiple-condition-met",
			state: WorldState{
				currentZoneId: testIdentity(t, "zone-one"),
				zones: []Zone{testZone(t, "zone-one")},
				flags: []Identity{testIdentity(t, "flag-one")},
				objects: []Object{
					testObject(t, "item-one", testInventoryZoneId),
					testObject(
						t, "object-one", "zone-one",
						testWithObjectProperty(t, "prop-one", "foo"),
						testWithObjectAction(VerbUse,
							testRule(t,
								[]Effect{SayEffect{"should-not-see-this"}},
								testWithConditionGroups(testConditionGroup(t,
									testWithPropertyCondition(t, "prop-one", "foo"),
									testWithFlagCondition(t, "flag-one", true),
								)),
							),
							testRule(t,
								[]Effect{SayEffect{"has-all-the-things"}},
								testWithConditionGroups(testConditionGroup(t,
									testWithUseWithCondition(t, "item-one"),
									testWithPropertyCondition(t, "prop-one", "foo"),
									testWithFlagCondition(t, "flag-one", true),
								)),
							),
						),
					),
				},
			},
			action: NewObjectActionWithItem(VerbUse, testIdentity(t, "object-one"), testIdentity(t, "item-one")),
			expected: []Effect{SayEffect{"has-all-the-things"}},
		},
	})
}
