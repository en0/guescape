package game

import (
	"testing"
)

func TestNewObject(t *testing.T) {

	object := NewObject(
		testIdentity(t, "obj-one"),
		testIdentity(t, "zone-one"),
		"utest", nil, nil, nil, nil, nil,
	)

	if !object.visible {
		t.Fatalf("Visibility defaults to %v, want %v", object.visible, true)
	}

	if !object.targetable {
		t.Fatalf("Targetability defaults to %v, want %v", object.targetable, true)
	}

	if object.carriable {
		t.Fatalf("Carriability defaults to %v, want %v", object.carriable, false)
	}
}
