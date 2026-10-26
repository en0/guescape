package hcl

import (
	"errors"
	"testing"

	"github.com/en0/guescape/internal/game"
)

type utilTestCase[T any] struct{
	given string
	expect T
	err error
}

func TestNewMessageKey(t *testing.T) {
	runUtilTests[game.MsgKey](t, NewMessageKey, []utilTestCase[game.MsgKey]{
		{ given: "CannotMove", 			expect: game.MsgCannotMove },
		{ given: "WayBlocked", 			expect: game.MsgWayBlocked },
		{ given: "DestinationMissing", 	expect: game.MsgDestinationMissing },
		{ given: "NoHandler", 			expect: game.MsgNoHandler },
		{ given: "ObjectNotFound", 		expect: game.MsgObjectNotFound },
		{ given: "ObjectNotTargetable",	expect: game.MsgObjectNotTargetable },
		{ given: "ObjectNotCarriable", 	expect: game.MsgObjectNotCarriable },
		{ given: "ObjectNotUsable", 	expect: game.MsgObjectNotUsable },
		{ given: "ObjectAlreadyTaken", 	expect: game.MsgObjectAlreadyTaken },
		{ given: "ObjectNotTaken", 		expect: game.MsgObjectNotTaken },
		{ given: "ItemNotFound", 		expect: game.MsgItemNotFound },
		{ given: "bad-val",				err: errBadMessageKey },
	})
}

func TestNewObjectVerb(t *testing.T) {
	runUtilTests[game.ObjectVerb](t, NewObjectVerb, []utilTestCase[game.ObjectVerb]{
		{ given: "examine", expect: game.VerbExamine },
		{ given: "take", 	expect: game.VerbTake },
		{ given: "use", 	expect: game.VerbUse },
		{ given: "drop", 	expect: game.VerbDrop },
		{ given: "bad-val",	err: errBadObjectVerb },
	})
}

func runUtilTests[T comparable](t *testing.T, fn func(string) (T, error), cases []utilTestCase[T]) {
	for _, tt := range cases {
		t.Run(tt.given, func(t *testing.T) {
			actual, err := fn(tt.given)
			if !errors.Is(err, tt.err) {
				t.Fatalf("unexpected error: %v, want %v", err, tt.err)
			}
			if actual != tt.expect {
				t.Fatalf("unexpected result: %v, want %v", actual, tt.expect)
			}
		})
	}
}
