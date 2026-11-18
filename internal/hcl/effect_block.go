package hcl

import (
	"errors"

	"github.com/en0/guescape/internal/game"
	"github.com/hashicorp/hcl/v2"
)

type effect interface {
	isEffect()
	toDomain() (game.Effect, error)
}

type effectBlock struct {
	Kind string   `hcl:"type,label"`
	Body hcl.Body `hcl:",remain"`
}

// sayEffectBlock

type sayEffectBlock struct {
	Text string `hcl:"text"`
}

func (sayEffectBlock) isEffect() { }

func (e sayEffectBlock) toDomain() (game.Effect, error) {
	return game.SayEffect{Text: e.Text}, nil
}

// currentZoneEffectBlock

type currentZoneEffectBlock struct {
	ZoneId string `hcl:"zone_id"`
}

func (currentZoneEffectBlock) isEffect() { }

func (e currentZoneEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.ZoneId)
	return game.CurrentZoneEffect{ZoneId: id}, err
}

// addFlagEffectBlock

type addFlagEffectBlock struct {
	FlagId 		string 	`hcl:"flag_id"`
}

func (addFlagEffectBlock) isEffect() { }

func (e addFlagEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.FlagId)
	return game.AddFlagEffect{Flag: id}, err
}

// removeFlagEffectBlock

type removeFlagEffectBlock struct {
	FlagId 		string 	`hcl:"flag_id"`
}

func (removeFlagEffectBlock) isEffect() { }

func (e removeFlagEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.FlagId)
	return game.RemoveFlagEffect{Flag: id}, err
}

// objectStateEffectBlock

type objectStateEffectBlock struct {
	ObjectId	string	`hcl:"object_id"`
	StateId 	string	`hcl:"state_id"`
	Value      	string	`hcl:"value"`
}

func (objectStateEffectBlock) isEffect() { }

func (e objectStateEffectBlock) toDomain() (game.Effect, error) {
	var errs []error
	oid, err := game.NewIdentity(e.ObjectId)
	if err != nil {errs = append(errs, err)}
	pid, err := game.NewIdentity(e.StateId)
	if err != nil {errs = append(errs, err)}
	return game.SetObjectStateEffect{
		ObjectId: oid,
		StateId: pid,
		Value: e.Value,
	}, errors.Join(errs...)
}

// objectVisibleEffectBlock

type objectVisibleEffectBlock struct {
	ObjectId	string	`hcl:"object_id"`
	Visible		bool	`hcl:"value"`
}

func (objectVisibleEffectBlock) isEffect() { }

func (e objectVisibleEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.ObjectId)
	return game.ObjectVisibleEffect{
		ObjectId: id,
		Visible: e.Visible,
	}, err
}

// objectTargetableEffectBlock

type objectTargetableEffectBlock struct {
	ObjectId 	string	`hcl:"object_id"`
	Targetable	bool	`hcl:"value"`
}

func (objectTargetableEffectBlock) isEffect() { }

func (e objectTargetableEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.ObjectId)
	return game.ObjectTargetableEffect{
		ObjectId: id,
		Targetable: e.Targetable,
	}, err
}

// objectCarriableEffectBlock

type objectCarriableEffectBlock struct {
	ObjectId 	string	`hcl:"object_id"`
	Carriable	bool	`hcl:"value"`
}

func (objectCarriableEffectBlock) isEffect() { }

func (e objectCarriableEffectBlock) toDomain() (game.Effect, error) {
	id, err := game.NewIdentity(e.ObjectId)
	return game.ObjectCarriableEffect{
		ObjectId: id,
		Carriable: e.Carriable,
	}, err
}

// objectLocationEffectBlock

type objectLocationEffectBlock struct {
	ObjectId 	string	`hcl:"object_id"`
	ZoneId 		string	`hcl:"zone_id"`
}

func (objectLocationEffectBlock) isEffect() {}

func (e objectLocationEffectBlock) toDomain() (game.Effect, error) {
	var errs []error
	oid, err := game.NewIdentity(e.ObjectId)
	if err != nil {errs = append(errs, err)}
	zid, err := game.NewIdentity(e.ZoneId)
	if err != nil {errs = append(errs, err)}
	return game.ObjectLocationEffect{
		ObjectId: oid,
		ZoneId: zid,
	}, errors.Join(errs...)
}

// endGameEffectBlock

type endGameEffectBlock struct {
	IsWin 	bool	`hcl:"is_win"`
	Text 	string	`hcl:"text"`
}

func (endGameEffectBlock) isEffect() { }

func (e endGameEffectBlock) toDomain() (game.Effect, error) {
	return game.EndGameEffect{
		IsWin: e.IsWin,
		Text: e.Text,
	}, nil
}

// activateModeEffectBlock

type activateModeEffectBlock struct {
	ModeId 	string 	`hcl:"mode_id"`
}

func (activateModeEffectBlock) isEffect() {}

func (e activateModeEffectBlock) toDomain() (game.Effect, error){
	id, err := game.NewIdentity(e.ModeId)
	return game.ActivateModeEffect{ModeId: id}, err
}

// clearModeEffectBlock

type clearModeEffectBlock struct { }

func (clearModeEffectBlock) isEffect() {}

func (e clearModeEffectBlock) toDomain() (game.Effect, error){
	return game.ClearModeEffect{}, nil
}

// setModeStateEffectBlock

type setModeStateEffectBlock struct {
	ModeId 	string	`hcl:"mode_id"`
	State 	string	`hcl:"state"`
	Value 	string	`hcl:"value"`
}

func (setModeStateEffectBlock) isEffect() {}

func (e setModeStateEffectBlock) toDomain() (game.Effect, error) {
	var errs []error
	mid, err := game.NewIdentity(e.ModeId)
	if err != nil {
		errs = append(errs, err)
	}
	sid, err := game.NewIdentity(e.State)
	if err != nil {
		errs = append(errs, err)
	}
	return game.SetModeStateEffect{
		ModeId: mid,
		State: sid,
		Value: e.Value,
	}, errors.Join(errs...)
}

// SetModePromptEffectBlock

type SetModePromptEffectBlock struct {
	ModeId 	string	`hcl:"mode_id"`
	Prompt 	string	`hcl:"prompt"`
}

func (SetModePromptEffectBlock) isEffect() {}

func (e SetModePromptEffectBlock) toDomain() (game.Effect, error) {
	var errs []error
	mid, err := game.NewIdentity(e.ModeId)
	if err != nil {
		errs = append(errs, err)
	}
	return game.SetModePromptEffect{
		ModeId: mid,
		Value: e.Prompt,
	}, errors.Join(errs...)
}
