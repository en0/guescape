package game

type MsgKey int

const (
	// Navigation Errors

	// MsgCannotMove reports that the chosen direction has no exit.
	MsgCannotMove          MsgKey = iota
	// MsgWayBlocked reports that an exit is blocked.
	MsgWayBlocked
	// MsgDestinationMissing reports that an exit's destination zone is missing.
	MsgDestinationMissing
	// MsgNoHandler reports that an action has no handler.
	MsgNoHandler

	// Object Errors

	// MsgObjectNotFound reports that the target object is not present.
	MsgObjectNotFound
	// MsgObjectNotTargetable reports that the object cannot be targeted.
	MsgObjectNotTargetable
	// MsgObjectNotCarriable reports that the object cannot be carried.
	MsgObjectNotCarriable
	// MsgObjectNotUsable reports that the object cannot be used.
	MsgObjectNotUsable
	// MsgObjectAlreadyTaken reports that the object is already in the inventory.
	MsgObjectAlreadyTaken
	// MsgObjectNotTaken reports that the object is not in the inventory.
	MsgObjectNotTaken

	// Item Errors

	// MsgItemNotFound reports that the item is not in the inventory.
	MsgItemNotFound
)

const (
	// Fallback message
	msgNoDefaultMessage = "The words for this moment seem to have gone missing"
)

var msgDefaults = map[MsgKey]string {
	// Navigation Error Defaults
	MsgCannotMove: "Cannot move that way",
	MsgWayBlocked: "That way is blocked",
	MsgDestinationMissing: "The way ahead blurs, like a page that was never written. You stay where you are.",
	MsgNoHandler: "That does not work here",

	// Object Errors Defaults
	MsgObjectNotFound: "Nothing like that here",
	MsgObjectNotTargetable: "Nothing more to make of it",
	MsgObjectNotCarriable: "Cannot carry that",
	MsgObjectNotUsable: "Cannot use that",
	MsgObjectAlreadyTaken: "Already in inventory",
	MsgObjectNotTaken: "Not in inventory",

	// Item Errors
	MsgItemNotFound: "You do not have that",
}

func message(key MsgKey) string {
	if m, ok := msgDefaults[key]; ok {
		return m
	}
	return msgNoDefaultMessage
}
