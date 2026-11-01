package lua

import (
	"testing"

	"github.com/en0/guescape/internal/game"
)

func testIdentity(t *testing.T, id string) game.Identity {
	t.Helper()
	v, err := game.NewIdentity(id)
	if err != nil {
		t.Fatalf("testIdentity: %v", err)
	}
	return v
}

func testIdentities(t *testing.T, names ...string) []game.Identity {
	t.Helper()
	var result []game.Identity
	for _, name := range names {
		id, err := game.NewIdentity(name)
		if err != nil {
			t.Fatalf("NewIdentity(%q): %v", name, err)
		}
		result = append(result, id)
	}
	return result
}

func testMockScriptContextFromContext(c *MockScriptContext) (*MockScriptContext) {
	if c != nil {
		return testMockScriptContext(
			c.hasFlag,
			c.propertyString,
			c.propertyBool,
			c.stateString,
			c.stateBool,
			c.currentModeId,
		)
	}
	return &MockScriptContext{}
}

func testMockScriptContext(
	hasFlag bool,
	propertyString string,
	propertyBool bool,
	stateString string,
	stateBool bool,
	currentModeId game.Identity,
) (*MockScriptContext) {
	return &MockScriptContext{
		hasFlag: hasFlag,
		propertyString: propertyString,
		propertyBool: propertyBool,
		stateString: stateString,
		stateBool: stateBool,
		currentModeId: currentModeId,
	}
}

type MockScriptContext struct {
	hasFlag bool
	propertyString string
	propertyBool bool
	stateString string
	stateBool bool
	currentModeId game.Identity
	hasFlagCalled bool
	hasFlagCalledWith game.Identity
	propertyCalled bool
	propertyCalledWith []game.Identity
	stateCalled bool
	stateCalledWith []game.Identity
}

func (m *MockScriptContext) HasFlag(flagId game.Identity) bool {
	m.hasFlagCalled = true
	m.hasFlagCalledWith = flagId
	return m.hasFlag
}

func (m *MockScriptContext) ObjectProperty(objId game.Identity, propId game.Identity) (string, bool) {
	m.propertyCalled = true
	m.propertyCalledWith = []game.Identity{objId, propId}
	return m.propertyString, m.propertyBool
}

func (m *MockScriptContext) ModeState(modeId game.Identity, state game.Identity) (string, bool) {
	m.stateCalled = true
	m.stateCalledWith = []game.Identity{modeId, state}
	return m.stateString, m.stateBool
}

func (m *MockScriptContext) CurrentModeId() game.Identity {
	return m.currentModeId
}
