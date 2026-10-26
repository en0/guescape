package game

// Effect is a change to apply to the world state.
type Effect interface {
	isEffect()
}

// SayEffect displays text to the player.
type SayEffect struct {
	// Text is the message to display.
	Text string
}
func (SayEffect) isEffect() {}

// CurrentZoneEffect moves the player to a zone.
type CurrentZoneEffect struct {
	// ZoneId is the destination zone.
	ZoneId Identity
}
func (CurrentZoneEffect) isEffect() {}

// AddFlagEffect sets a world flag.
type AddFlagEffect struct {
	// Flag is the flag to set.
	Flag Identity
}
func (AddFlagEffect) isEffect() {}

// RemoveFlagEffect clears a world flag.
type RemoveFlagEffect struct {
	// Flag is the flag to clear.
	Flag Identity
}
func (RemoveFlagEffect) isEffect() {}

// ObjectPropertyEffect sets a property on an object.
type ObjectPropertyEffect struct {
	// ObjectId is the target object.
	ObjectId 	Identity
	// PropertyId is the property to set.
	PropertyId 	Identity
	// Value is the value to set.
	Value 		string
}
func (ObjectPropertyEffect) isEffect() {}

// ObjectVisibleEffect sets whether an object is visible.
type ObjectVisibleEffect struct {
	// ObjectId is the target object.
	ObjectId 	Identity
	// Visible is the visibility to set.
	Visible 	bool
}
func (ObjectVisibleEffect) isEffect() {}

// ObjectTargetableEffect sets whether an object can be targeted.
type ObjectTargetableEffect struct {
	// ObjectId is the target object.
	ObjectId 	Identity
	// Targetable is the targetability to set.
	Targetable 	bool
}
func (ObjectTargetableEffect) isEffect() {}

// TODO: Do we need this?
// type ObjectCarriableEffect struct {
// 	ObjectId Identity
// 	Carriable bool
// }
// func (ObjectCarriableEffect) isEffect() {}

// ObjectLocationEffect moves an object to a zone.
type ObjectLocationEffect struct {
	// ObjectId is the object to move.
	ObjectId 	Identity
	// ZoneId is the destination zone. Can be InventoryZone
	ZoneId 		Identity
}
func (ObjectLocationEffect) isEffect() {}

// EndGameEffect ends the game.
type EndGameEffect struct {
	// IsWin reports whether the game was won.
    IsWin 	bool
	// Text is the message to display.
    Text 	string
}
func (EndGameEffect) isEffect() {}
