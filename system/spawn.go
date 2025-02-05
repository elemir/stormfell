package system

import (
	"errors"
	"fmt"
	"image"
	"iter"

	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/gloomo/input"
	gmodel "github.com/elemir/gloomo/model"

	"github.com/elemir/stormfell/model"
)

var errAnimationNotLoaded = errors.New("animation is not loaded")

type MouseInput interface {
	IsPressed(button input.MouseButton) bool
	Position() image.Point
}

type AnimationLoader interface {
	Load(path string) (*gmodel.Animation, bool)
}

type UnitRepo interface {
	Upsert(id gid.ID, unit model.Unit)
	List() iter.Seq2[gid.ID, model.Unit]
}

type IDGenerator interface {
	New() gid.ID
}

type SpawnWarrior struct {
	IDGen           IDGenerator
	MouseInput      MouseInput
	UnitRepo        UnitRepo
	AnimationLoader AnimationLoader
}

func (sw *SpawnWarrior) Run() error {
	// TODO(elemir): loaders should return something that can be persist in all cases
	// TODO(elemir): we should have special unit spec, not animation itself
	unitAnim, ok := sw.AnimationLoader.Load("unit.yaml")
	if !ok {
		return fmt.Errorf("animation %q: %w", "unit.yaml", errAnimationNotLoaded)
	}

	if !sw.MouseInput.IsPressed(input.MouseButtonLeft) {
		return nil
	}

	id := sw.IDGen.New()
	sw.UnitRepo.Upsert(id, model.Unit{
		Animation: unitAnim,
		Position:  sw.MouseInput.Position().Sub(unitAnim.Size.Div(2)),
	})

	return nil
}
