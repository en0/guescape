package game

// PropertyCondition matches an object property against an expected value.
type PropertyCondition struct {
	id  	Identity
	value 	string
}

// NewPropertyCondition returns a PropertyCondition for the given property and value.
func NewPropertyCondition(id Identity, value string) PropertyCondition {
	return PropertyCondition{id: id, value: value}
}

// FlagCondition matches if a flag is present in the world.
type FlagCondition struct {
	id		Identity
	value 	bool
}

// NewFlagCondition returns a FlagCondition for the given flag and value.
func NewFlagCondition(id Identity, value bool) FlagCondition {
	return FlagCondition{id: id, value: value}
}

// ConditionGroup matches when useWith and all flag and property conditions match.
//A nil useWith matches when no item is used.
type ConditionGroup struct {
	properties []PropertyCondition // All must match
	flags      []FlagCondition     // All must match
	useWith    Identity            // Must match, if specified
}

// NewConditionGroup returns a condition group for the given useWith, flag, and proerty conditions
func NewConditionGroup(
	useWith 			Identity,
	flagConditions 		[]FlagCondition,
	propertyConditions 	[]PropertyCondition,
) ConditionGroup {
	return ConditionGroup{
		useWith: 	useWith,
		flags: 		flagConditions,
		properties: propertyConditions,
	}
}

func (cg ConditionGroup) matches(ws WorldState, obj Object, itemId Identity) bool {
	allProperties := existsAll(cg.properties, func(pc PropertyCondition) bool {
		p, ok := obj.property(pc.id)
		return ok && p == pc.value
	})
	allFlags := existsAll(cg.flags, func(fc FlagCondition) bool {
		return ws.hasFlag(fc.id) == fc.value
	})
	return cg.useWith == itemId && allProperties && allFlags
}

