package game

// InventoryZoneId is the zone that holds objects carried by the player.
var InventoryZoneId = Identity{value: "player-inventory"}

// Zone is a location in the world with exits to other zones.
type Zone struct {
	id   Identity
	desc string
	exits map[Identity]ZoneExit
}

// NewZone returns a Zone with the given id, description, and exits.
func NewZone(id Identity, desc string, exits map[Identity]ZoneExit) Zone {
	return Zone{id: id, desc: desc, exits: exits}
}

// ZoneExit is a path to a destination zone, optionally gated by the given flags.
type ZoneExit struct {
	zoneId Identity
	flags []Identity
}

// NewZoneExit returns a ZoneExit to zoneId gated by the given flags.
func NewZoneExit(zoneId Identity, flags []Identity) ZoneExit {
	return ZoneExit{zoneId: zoneId, flags: flags}
}
