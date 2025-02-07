package repo

import (
	"iter"
	"math"

	"github.com/elemir/gloomo/geom"
	gid "github.com/elemir/gloomo/id"
	gmodel "github.com/elemir/gloomo/model"

	"github.com/elemir/stormfell/model"
)

type Collection[T any] interface {
	Set(id gid.ID, val T)
	Get(id gid.ID) (T, bool)
	Items() iter.Seq2[gid.ID, T]
}

type Unit struct {
	ZIndices          Collection[int]
	CurrentAnimations Collection[string]
	Animations        Collection[*gmodel.AnimationSheet]
	Positions         Collection[geom.Vec2]
	Velocities        Collection[geom.Vec2]
}

func (u *Unit) List() iter.Seq2[gid.ID, model.Unit] {
	return func(yield func(gid.ID, model.Unit) bool) {
		for id, pos := range u.Positions.Items() {
			anim, ok := u.Animations.Get(id)
			if !ok {
				continue
			}

			vel, ok := u.Velocities.Get(id)
			if !ok {
				continue
			}

			unit := model.Unit{
				Position:  pos,
				Velocity:  vel,
				Animation: anim,
			}

			if !yield(id, unit) {
				return
			}
		}
	}
}

var animations = []string{
	"north",
	"northeast",
	"east",
	"southeast",
	"south",
	"southwest",
	"west",
	"northwest",
}

func (u *Unit) Upsert(id gid.ID, unit model.Unit) {
	dir := direction(unit.Velocity.Angle())

	u.Velocities.Set(id, unit.Velocity)
	u.Positions.Set(id, unit.Position)
	u.CurrentAnimations.Set(id, animations[dir])
	u.Animations.Set(id, unit.Animation)
	u.ZIndices.Set(id, 1)
}

func (u *Unit) Get(id gid.ID) (model.Unit, bool) {
	pos, ok := u.Positions.Get(id)
	if !ok {
		return model.Unit{}, false
	}

	anim, ok := u.Animations.Get(id)
	if !ok {
		return model.Unit{}, false
	}

	vel, ok := u.Velocities.Get(id)
	if !ok {
		return model.Unit{}, false
	}

	return model.Unit{
		Position:  pos,
		Animation: anim,
		Velocity:  vel,
	}, true
}

func direction(angle geom.Angle) int {
	angle = angle.Normalize()

	direction := int(math.Round(float64(angle*4/math.Pi))) + 4
	if direction == 8 {
		direction = 0
	}

	return direction
}
