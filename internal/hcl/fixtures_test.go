package hcl

import "testing"

type testOptions[T any] func(*T)

func applyTestOptions[T any](v T, opts ...testOptions[T]) T {
	for _, opt := range opts { opt(&v) }
	return v
}

func withWorldZones(zones ...zoneBlock) testOptions[worldBlock] {
	return func(wb *worldBlock) {wb.Zones = append(wb.Zones, zones...)}
}

func withWorldObjects(objs ...objectBlock) testOptions[worldBlock] {
	return func(wb *worldBlock) {wb.Objects = append(wb.Objects, objs...)}
}

func testWorldBlock(t *testing.T, startZoneId string, desc string, opts ...testOptions[worldBlock]) worldBlock {
	t.Helper()
	return applyTestOptions(worldBlock{
		StartZoneId: startZoneId,
		Desc: desc,
	}, opts...)
}

func withExitFlag(flag string) testOptions[exitBlock] {
	return func(eb *exitBlock) {
		eb.Flags = append(eb.Flags, flag)
	}
}

func withExitFlags(flags ...string) testOptions[exitBlock] {
	return func(eb *exitBlock) {
		eb.Flags = flags
	}
}

func testExitBlock(t *testing.T, id string, to string, opts ...testOptions[exitBlock]) exitBlock {
	t.Helper()
	return applyTestOptions(exitBlock{
		Id: id,
		ZoneId: to,
	}, opts...)
}

func withZoneExit(exit exitBlock) testOptions[zoneBlock] {
	return func(zb *zoneBlock) {
		zb.Exits = append(zb.Exits, exit)
	}
}

func testZoneBlock(t *testing.T, id string, opts ...testOptions[zoneBlock]) zoneBlock {
	t.Helper()
	return applyTestOptions(zoneBlock{
		Id: id,
		Desc: id,
	}, opts...)
}

func withObjectDesc(desc string) testOptions[objectBlock] {
	return func(ob *objectBlock) {ob.Desc = &desc}
}

func withObjectVisible(value bool) testOptions[objectBlock] {
	return func(ob *objectBlock) {ob.Visible = &value}
}

func withObjectTargetable(value bool) testOptions[objectBlock] {
	return func(ob *objectBlock) {ob.Targetable = &value}
}

func withObjectCarriable(value bool) testOptions[objectBlock] {
	return func(ob *objectBlock) {ob.Carriable = &value}
}

func withObjectProperty(id string, value string) testOptions[objectBlock] {
	return func(ob *objectBlock) {
		ob.Properties = append(ob.Properties, propertyBlock{
			Key: id,
			Value: value,
		})
	}
}

func withObjectActions(actions ...actionBlock) testOptions[objectBlock] {
	return func(ob *objectBlock) {
		ob.Actions = append(ob.Actions, actions...)
	}
}

func testObjectBlock(t *testing.T, id string, zoneId string, opts ...testOptions[objectBlock]) objectBlock {
	t.Helper()
	return applyTestOptions(objectBlock{
		Id: id,
		ZoneId: zoneId,
	}, opts...)
}

func testActionBlock(t *testing.T, verb string, rules ...ruleBlock) actionBlock {
	t.Helper()
	return applyTestOptions(actionBlock{Verb: verb, Rules: rules})
}

func withRuleConditions(t *testing.T, condition ...conditionGroupBlock) testOptions[ruleBlock] {
	return func(rb *ruleBlock) {
		rb.ConditionGroups = condition
	}
}

func withRuleEffects(t *testing.T, effects ...effect) testOptions[ruleBlock] {
	return func(rb *ruleBlock) {
		rb.Effects = effects
	}
}

func testRuleBlock(t *testing.T, opts ...testOptions[ruleBlock]) ruleBlock {
	t.Helper()
	return applyTestOptions(ruleBlock{}, opts...)
}

func withFlagCondition(id string, value bool) testOptions[conditionGroupBlock] {
	return func(wb *conditionGroupBlock) {
		wb.FlagConditions = append(wb.FlagConditions, flagConditionBlock{
			Id: id,
			Value: value,
		})
	}
}

func withPropCondition(id string, value string) testOptions[conditionGroupBlock] {
	return func(wb *conditionGroupBlock) {
		wb.PropertyConditions = append(wb.PropertyConditions, propertyConditionBlock{
			Id: id,
			Value: value,
		})
	}
}

func withUseCondition(item string) testOptions[conditionGroupBlock] {
	return func(wb *conditionGroupBlock) {wb.UseWith=&item}
}

func testWhenBlock(t *testing.T, opts ...testOptions[conditionGroupBlock]) conditionGroupBlock {
	t.Helper()
	return applyTestOptions(conditionGroupBlock{}, opts...)
}
