package game

import (
	"reflect"
	"testing"
)

func TestNewNavActions(t *testing.T) {
	identOne := testIdentity(t, "ident-one")
	expect := Action{
		kind: navActionType,
		navCommand: CmdMove,
		exitId: identOne,
	}
	actual := NewNavAction(CmdMove, identOne)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("NewNavAction(%v, %q) = %v, want %v", CmdMove, identOne, actual, expect)
	}
}

func TestNewObjectAction(t *testing.T) {
	identOne := testIdentity(t, "ident-one")
	expect := Action{
		kind: objectActionType,
		objectVerb: VerbExamine,
		objectTarget: identOne,
		objectHasItem: false,
	}
	actual := NewObjectAction(VerbExamine, identOne)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("NewObjectAction(%v, %q) = %v, want %v", CmdMove, identOne, actual, expect)
	}
}

func TestNewObjectWithItemAction(t *testing.T) {
	identOne := testIdentity(t, "ident-one")
	identTwo := testIdentity(t, "ident-two")
	expect := Action{
		kind: objectActionType,
		objectVerb: VerbExamine,
		objectTarget: identOne,
		objectItem: identTwo,
		objectHasItem: true,
	}
	actual := NewObjectActionWithItem(VerbExamine, identOne, identTwo)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("NewObjectActionWithItem(%v, %q, %q) = %v, want %v", CmdMove, identOne, identTwo, actual, expect)
	}
}
