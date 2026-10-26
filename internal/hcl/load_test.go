package hcl

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
)

func mustParse(t *testing.T, src []byte) hcl.Body {
	t.Helper()
	body, err := parseHCL(src, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestDecodeWorld(t *testing.T) {

	exitToZoneTwo := testExitBlock(t, "exit-to-zone-two", "zone-two", withExitFlags("flag-one", "flag-two"))
	tests := []struct{
		name string
		hcl []byte
		expect worldBlock
	} {
		{
			name: "minimal-world",
			hcl: []byte(`game {
				start_zone_id = "zone-one"
				desc = "world-one"
				zone "zone-one" {
					desc = "zone-one"
					exit "exit-to-zone-two" {
						flags = ["flag-one", "flag-two"]
						to = "zone-two"
					}
				}
				zone "zone-two" {
					desc = "zone-two"
				}
				object "object-one" {
					zone_id = "zone-one"
					desc = "object-one"
					visible = false
					targetable = false
					carriable = true
					property "prop-one" {
						value = "prop-one"
					}
					action "examine" {
						rule {
							when {
								property "prop-one" {
									value = "prop-one"
								}
								flag "flag-one" {
									value = true
								}
								use_with = "item-one"
							}
							effect "say" {
								text = "say-one"
							}
						}
					}
				}
			}`),
			expect: testWorldBlock(t, "zone-one", "world-one",
				withWorldZones(
					testZoneBlock(t, "zone-one", withZoneExit(exitToZoneTwo)),
					testZoneBlock(t, "zone-two"),
				),
				withWorldObjects(
					testObjectBlock(t, "object-one", "zone-one",
						withObjectVisible(false),
						withObjectTargetable(false),
						withObjectCarriable(true),
						withObjectDesc("object-one"),
						withObjectProperty("prop-one", "prop-one"),
						withObjectActions(
							testActionBlock(t, "examine",
								testRuleBlock(t,
									withRuleEffects(t, sayEffectBlock{"say-one"}),
									withRuleConditions(t,
										testWhenBlock(t,
											withPropCondition("prop-one", "prop-one"),
											withFlagCondition("flag-one", true),
											withUseCondition("item-one"),
										),
									),
								),
							),
						),
					),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := decodeWorld(mustParse(t, tt.hcl))
			if err != nil {t.Fatalf("%v", err)}
			if !reflect.DeepEqual(actual, tt.expect) {
				t.Fatalf("LoadWorld() -> %v, want %v", actual, tt.expect)
			}
		})
	}
}

// func TestToWorldState(t *testing.T) {
// 	zoneOneId, err := game.NewIdentity("zone-one")
// 	if err != nil { t.Fatal(err) }
// 	w, err := game.NewWorldState(zoneOneId, "world-one", nil, nil)
// 	if err != nil { t.Fatal(err) }
// 	tests := []struct {
// 		name string
// 		input worldBlock
// 		expect game.WorldState
// 	} {
// 		{
// 			name: "world-only",
// 			input: testWorldBlock(t, "zone-one", "world-one"),
// 			expect: w,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			actual, err := toWorldState(tt.input)
// 			if err != nil {t.Fatalf("%v", err)}
// 			if !reflect.DeepEqual(actual, tt.expect) {
// 				t.Fatalf("toWorldState() -> %v, want %v", actual, tt.expect)
// 			}
// 		})
// 	}
// }
