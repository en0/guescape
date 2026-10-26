package game

import "testing"

func TestConditionGroupMatches(t *testing.T) {
	tests := []struct{
		name string
		expect bool
		cg ConditionGroup
		ws WorldState
		obj Object
		item Identity
	} {
		{
			name: "match-flag",
			expect: true,
			cg: testConditionGroup(t,
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(
				t, "zone-one", "desc", []Zone{
					testZone(t, "zone-one"),
				}, nil, testWithWorldStateFlag(t, "flag-one"),
			),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "not-match-flag",
			expect: false,
			cg: testConditionGroup(t,
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "match-not-flag",
			expect: true,
			cg: testConditionGroup(t,
				testWithFlagCondition(t, "flag-one", false),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "not-match-not-flag",
			expect: false,
			cg: testConditionGroup(t,
				testWithFlagCondition(t, "flag-one", false),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{
					testZone(t, "zone-one"),
				}, nil, testWithWorldStateFlag(t, "flag-one"),
			),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "match-prop",
			expect: true,
			cg: testConditionGroup(t,
				testWithPropertyCondition(t, "prop-one", "foo"),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
		{
			name: "not-match-prop",
			expect: false,
			cg: testConditionGroup(t,
				testWithPropertyCondition(t, "prop-one", "bar"),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
		{
			name: "match-item",
			expect: true,
			item: testIdentity(t, "item-one"),
			cg: testConditionGroup(t, testWithUseWithCondition(t, "item-one")),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "not-match-item",
			expect: false,
			cg: testConditionGroup(t, testWithUseWithCondition(t, "item-one")),
			ws: testWorldState(t, "zone-one", "desc", []Zone{testZone(t, "zone-one")}, nil),
			obj: testObject(t, "obj-one", "zone-one"),
		},
		{
			name: "match-all",
			expect: true,
			item: testIdentity(t, "item-one"),
			cg: testConditionGroup(t,
				testWithUseWithCondition(t, "item-one"),
				testWithPropertyCondition(t, "prop-one", "foo"),
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{
				testZone(t, "zone-one"),
			}, nil, testWithWorldStateFlag(t, "flag-one")),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
		{
			name: "not-match-all-missing-item",
			expect: false,
			cg: testConditionGroup(t,
				testWithUseWithCondition(t, "item-one"),
				testWithPropertyCondition(t, "prop-one", "foo"),
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{
				testZone(t, "zone-one"),
			}, nil, testWithWorldStateFlag(t, "flag-one")),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
		{
			name: "not-match-all-missing-flag",
			expect: false,
			item: testIdentity(t, "item-one"),
			cg: testConditionGroup(t,
				testWithUseWithCondition(t, "item-one"),
				testWithPropertyCondition(t, "prop-one", "foo"),
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{
				testZone(t, "zone-one"),
			}, nil),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
		{
			name: "not-match-all-wrong-prop-value",
			expect: false,
			item: testIdentity(t, "item-one"),
			cg: testConditionGroup(t,
				testWithUseWithCondition(t, "item-one"),
				testWithPropertyCondition(t, "prop-one", "bar"),
				testWithFlagCondition(t, "flag-one", true),
			),
			ws: testWorldState(t, "zone-one", "desc", []Zone{
				testZone(t, "zone-one"),
			}, nil, testWithWorldStateFlag(t, "flag-one")),
			obj: testObject(t, "obj-one", "zone-one",
				testWithObjectProperty(t, "prop-one", "foo"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.cg.matches(tt.ws, tt.obj, tt.item)
			if tt.expect != actual {
				t.Fatalf("ConditionGroup.match() = %v, wants %v", actual, tt.expect)
			}
		})
	}
}
