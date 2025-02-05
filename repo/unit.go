package repo

import (
	"image"
	"iter"

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
	Animations        Collection[*gmodel.Animation]
	Positions         Collection[image.Point]
}

func (u *Unit) List() iter.Seq2[gid.ID, model.Unit] {
	return func(yield func(gid.ID, model.Unit) bool) {
		for id, pos := range u.Positions.Items() {
			anim, ok := u.Animations.Get(id)
			if !ok {
				continue
			}

			unit := model.Unit{
				Position:  pos,
				Animation: anim,
			}

			if !yield(id, unit) {
				return
			}
		}
	}
}

func (u *Unit) Upsert(id gid.ID, unit model.Unit) {
	u.Positions.Set(id, unit.Position)
	u.CurrentAnimations.Set(id, "south")
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

	return model.Unit{
		Position:  pos,
		Animation: anim,
	}, true
}
