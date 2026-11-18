package game

// Object is a thing in the world that the player can interact with.
type Object struct {
	id         Identity
	zoneId     Identity
	desc       string
	visible    bool
	targetable bool
	carriable  bool
	states map[Identity]string
	actions    map[ObjectVerb][]Rule
}

func (o Object) inInventory() bool {
	return o.zoneId == InventoryZoneId
}

func (o Object) state(id Identity) (string, bool) {
	result, ok := o.states[id]
	return result, ok
}

func (o Object) hasState(id Identity) bool {
	_, ok := o.state(id)
	return ok
}

// NewObject returns an Object. A nil visible, targetable, or carriable defaults
// to true, true, and false respectively.
func NewObject(
	id 			Identity,
	zoneId 		Identity,
	desc 		string,
	visible 	*bool,
	targetable 	*bool,
	carriable 	*bool,
	states 		map[Identity]string,
	actions 	map[ObjectVerb][]Rule,
) (Object){
	ptr := func (b bool) *bool {return &b}
	if visible == nil {visible = ptr(true)}
	if targetable == nil {targetable = ptr(true)}
	if carriable == nil {carriable = ptr(false)}
	return Object{
		id: 		id,
		zoneId: 	zoneId,
		desc: 		desc,
		visible: 	*visible,
		targetable: *targetable,
		carriable: 	*carriable,
		states: states,
		actions: 	actions,
	}
}
