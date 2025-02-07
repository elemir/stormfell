package system

import (
	"errors"
	"fmt"
	"image"
	"iter"
	"math/rand/v2"

	"github.com/elemir/gloomo/geom"
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
	Load(path string) (*gmodel.AnimationSheet, bool)
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

	// TODO(elemir): remove me, only for debugging
	vel := image.Pt(rand.IntN(3)-1, rand.IntN(3)-1)

	id := sw.IDGen.New()
	sw.UnitRepo.Upsert(id, model.Unit{
		Animation: unitAnim,
		Position:  geom.FromPoint(sw.MouseInput.Position().Sub(unitAnim.Size.Div(2))),
		Velocity:  geom.FromPoint(vel),
	})

	return nil
}
