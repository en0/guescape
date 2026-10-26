package game

// Rule pairs condition groups with the effects applied when ANY group matches.
type Rule struct {
 	// Any condition group may match
	conditionGroups []ConditionGroup
	effects 		[]Effect
}

// NewRule returns a Rule with the given condition groups and effects.
func NewRule(
	conditionGroups []ConditionGroup,
	effects 		[]Effect,
) Rule {
	return Rule{
		conditionGroups: 	conditionGroups,
		effects: 			effects,
	}
}

func (r Rule) effectsFor(ws WorldState, obj Object, itemId Identity) []Effect {
	if len(r.conditionGroups) == 0 || exists(r.conditionGroups, func(gc ConditionGroup) bool {
		return gc.matches(ws, obj, itemId)
	}) {
		return r.effects
	}
	return nil
}
