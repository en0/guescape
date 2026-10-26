package game

// PlanEffects returns the effects produced by applying act to ws.
func PlanEffects(ws WorldState, act Action) []Effect {
	switch act.kind {
	case navActionType:
		return planNavEffects(ws, act)
	case objectActionType:
		return planObjectEffects(ws, act)
	}
	return []Effect{}
}

func planNavEffects(ws WorldState, act Action) []Effect {
	switch act.navCommand {
	case CmdMove:
		return planMoveEffects(ws, act)
	}
	return []Effect{}
}

func planObjectEffects(ws WorldState, act Action) []Effect {

	obj, found := ws.actionableObject(act.objectTarget)

	if !found {
		return sayMessage(ws, MsgObjectNotFound)
	}

	if act.objectHasItem {
		item, itemFound := ws.object(act.objectItem)
		if !itemFound || !item.inInventory(){
			return sayMessage(ws, MsgItemNotFound)
		}
	}

	if rules, ok := obj.actions[act.objectVerb]; ok {
		return planCustomEffect(ws, obj, act.objectItem, rules)
	}

	// Only custom handlers can deal with 'on <item>' commands
	if act.objectHasItem {
		return sayMessage(ws, MsgNoHandler)
	}

	switch act.objectVerb {
	case VerbExamine:
		return []Effect{SayEffect{obj.desc}}
	case VerbTake:
		return planTakeEffects(ws, obj)
	case VerbDrop:
		return planDropEffects(ws, obj)
	case VerbUse:
		return sayMessage(ws, MsgObjectNotUsable)
	}

	return sayMessage(ws, MsgNoHandler)
}

func planMoveEffects(ws WorldState, act Action) []Effect {

	currentZone := ws.currentZone()
	exit, found := currentZone.exits[act.exitId]

	if !found {
	 	return sayMessage(ws, MsgCannotMove)
	}

	if !existsAll(exit.flags, func(flag Identity) bool {return ws.hasFlag(flag)}) {
		return sayMessage(ws, MsgWayBlocked)
	}

	if !ws.hasZone(exit.zoneId) {
		// Should never happen
		return sayMessage(ws, MsgDestinationMissing)
	}

	return []Effect{CurrentZoneEffect{ZoneId: exit.zoneId}}
}

func planTakeEffects(ws WorldState, obj Object) []Effect {
	if obj.zoneId == InventoryZoneId {
		return sayMessage(ws, MsgObjectAlreadyTaken)
	}
	if !obj.carriable {
		return sayMessage(ws, MsgObjectNotCarriable)
	}
	return []Effect{
		ObjectLocationEffect{
			ObjectId: obj.id,
			ZoneId: InventoryZoneId,
		},
	}
}

func planDropEffects(ws WorldState, obj Object) []Effect {
	if obj.zoneId != InventoryZoneId {
		return sayMessage(ws, MsgObjectNotTaken)
	}
	return []Effect{
		ObjectLocationEffect{
			ObjectId: obj.id,
			ZoneId: ws.currentZoneId,
		},
	}
}

func planCustomEffect(ws WorldState, obj Object, itemId Identity, rules []Rule) []Effect {
	var result []Effect
	for _, rule := range rules {
		result = append(result, rule.effectsFor(ws, obj, itemId)...)
	}
	if len(result) > 0 {
		return result
	}
	// This case means that the custom action didn't match any rules (NoHandler).
	return sayMessage(ws, MsgNoHandler)
}

func sayMessage(ws WorldState, key MsgKey) []Effect {
	return []Effect{SayEffect{ws.message(key)}}
}
