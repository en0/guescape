package game

// StateCondition matches an object state against an expected value.
type StateCondition struct {
	id  	Identity
	value 	string
}

// NewStateCondition returns a ObjectStateCondition for the given state and value.
func NewStateCondition(id Identity, value string) StateCondition {
	return StateCondition{id: id, value: value}
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

// ConditionGroup matches when useWith and all flag and state conditions match.
//A nil useWith matches when no item is used.
type ConditionGroup struct {
	states 	   []StateCondition // All must match
	flags      []FlagCondition     // All must match
	useWith    Identity            // Must match, if specified
}

// NewConditionGroup returns a condition group for the given useWith, flag, and proerty conditions
func NewConditionGroup(
	useWith 			Identity,
	flagConditions 		[]FlagCondition,
	stateConditions 	[]StateCondition,
) ConditionGroup {
	return ConditionGroup{
		useWith: 	useWith,
		flags: 		flagConditions,
		states: stateConditions,
	}
}

func (cg ConditionGroup) matches(ws WorldState, obj Object, itemId Identity) bool {
	allStates := existsAll(cg.states, func(sc StateCondition) bool {
		s, ok := obj.state(sc.id)
		return ok && s == sc.value
	})
	allFlags := existsAll(cg.flags, func(fc FlagCondition) bool {
		return ws.HasFlag(fc.id) == fc.value
	})
	return cg.useWith == itemId && allStates && allFlags
}

