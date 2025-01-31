package repo

import (
	"image"
	"iter"

	gid "github.com/elemir/gloomo/id"

	"github.com/elemir/stormfell/model"
)

type Collection[T any] interface {
	Set(id gid.ID, val T)
	Get(id gid.ID) (T, bool)
	Items() iter.Seq2[gid.ID, T]
}

type Unit struct {
	Positions Collection[image.Point]
}

func (u *Unit) List() iter.Seq2[gid.ID, model.Unit] {
	return func(yield func(gid.ID, model.Unit) bool) {
		for id, pos := range u.Positions.Items() {
			unit := model.Unit{
				Position: pos,
			}

			if !yield(id, unit) {
				return
			}
		}
	}
}

func (u *Unit) Upsert(id gid.ID, unit model.Unit) {
	u.Positions.Set(id, unit.Position)
}

func (u *Unit) Get(id gid.ID) (model.Unit, bool) {
	pos, ok := u.Positions.Get(id)
	if !ok {
		return model.Unit{}, false
	}

	return model.Unit{
		Position: pos,
	}, true
}
