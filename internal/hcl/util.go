package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
)

var errBadMessageKey = errors.New("bad message key")

func NewMessageKey(v string) (game.MsgKey, error) {
	switch v {
	case "CannotMove":
		return game.MsgCannotMove ,nil
	case "WayBlocked":
		return game.MsgWayBlocked ,nil
	case "DestinationMissing":
		return game.MsgDestinationMissing ,nil
	case "NoHandler":
		return game.MsgNoHandler ,nil
	case "ObjectNotFound":
		return game.MsgObjectNotFound ,nil
	case "ObjectNotTargetable":
		return game.MsgObjectNotTargetable ,nil
	case "ObjectNotCarriable":
		return game.MsgObjectNotCarriable ,nil
	case "ObjectNotUsable":
		return game.MsgObjectNotUsable ,nil
	case "ObjectAlreadyTaken":
		return game.MsgObjectAlreadyTaken ,nil
	case "ObjectNotTaken":
		return game.MsgObjectNotTaken ,nil
	case "ItemNotFound":
		return game.MsgItemNotFound ,nil
	}
	return game.MsgKey(0), errBadMessageKey
}

var errBadObjectVerb = errors.New("bad object verb")

func NewObjectVerb(v string) (game.ObjectVerb, error) {
	switch v {
	case "examine":
		return game.VerbExamine ,nil
	case "take":
		return game.VerbTake ,nil
	case "use":
		return game.VerbUse ,nil
	case "drop":
		return game.VerbDrop ,nil
	}
	return game.ObjectVerb(0), errBadObjectVerb
}
