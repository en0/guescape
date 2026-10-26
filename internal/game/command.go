package game

type actionType int

const (
	navActionType actionType = iota
	objectActionType
)

// NavCommand is a movement command issued by the player, such as CmdMove
type NavCommand int

// ObjectVerb is an action a player performs on an object, such as VerbUse
type ObjectVerb int

const (
	// CmdMove moves the player through an exit.
	CmdMove      NavCommand   = iota
)

const (
	// VerbExamine describes an object.
	VerbExamine  ObjectVerb   = iota
	// VerbTake moves an object into the inventory.
	VerbTake
	// VerbUse uses an object.
	VerbUse
	// VerbDrop moves an object out of the inventory.
	VerbDrop
)

// Action is a player command to resolve against the world state.
type Action struct {
	kind actionType

	// NavActionType fields
	navCommand NavCommand
	exitId Identity

	// ObjectActionType fields
	objectVerb ObjectVerb
	objectTarget Identity
	objectHasItem bool
	objectItem Identity
}

// NewNavAction returns a navigation Action for the given command and exit.
func NewNavAction(cmd NavCommand, exitId Identity) Action {
	return Action{kind: navActionType, navCommand: cmd, exitId: exitId}
}

// NewObjectAction returns an object Action for the given verb and object.
func NewObjectAction(verb ObjectVerb, objectId Identity) Action {
	return Action{
		kind: objectActionType,
		objectVerb: verb,
		objectTarget: objectId,
		objectHasItem: false,
	}
}

// NewObjectActionWithItem returns an object Action for the given verb with attached item.
func NewObjectActionWithItem(verb ObjectVerb, objectId Identity, itemId Identity) Action {
	return Action{
		kind: objectActionType,
		objectVerb: verb,
		objectTarget: objectId,
		objectItem: itemId,
		objectHasItem: true,
	}
}
