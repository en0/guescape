package game

// TODO: Maybe, add currentZone?

// The ScriptContext interface provides read-only access to WorldState
type ScriptContext interface {
	// Check if the palyer has the given flag
	HasFlag(Identity) bool
	// Retrieve, from the given object, the value of the given property.
	ObjectProperty(Identity, Identity) (string, bool)
	// Retrieve, from the given state value from the current mode.
	CurrentModeState(Identity) (string, bool)
}

// Port definition to execute a compiled script.
type CompiledScript interface {
	// Run the script with the given input.
	Run(ctx ScriptContext, args []string) ([]Effect, error)
}
