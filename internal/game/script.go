package game

// TODO: Maybe, add currentZone?

// The ScriptContext interface provides read-only access to WorldState
type ScriptContext interface {
	// Check if the player has the given flag
	HasFlag(Identity) bool
	// Retrieve, from the given object, the value of the given property.
	ObjectProperty(Identity, Identity) (string, bool)
	// Retrieve, from the given mode, the value of the given state.
	ModeState(Identity, Identity) (string, bool)
	// Retrieve the id of the active mode
	CurrentModeId() (Identity)
}

// Port definition to execute a compiled script.
type CompiledScript interface {
	// Run the script with the given input.
	Run(ctx ScriptContext, args []string) ([]Effect, error)
}
