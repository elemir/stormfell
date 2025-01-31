package system

import (
	"image"

	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/gloomo/input"

	"github.com/elemir/stormfell/model"
)

type MouseInput interface {
	IsPressed(button input.MouseButton) bool
	Position() image.Point
}

type UnitRepo interface {
	Upsert(id gid.ID, unit model.Unit)
}

type IDGenerator interface {
	New() gid.ID
}

type SpawnWarrior struct {
	IDGen      IDGenerator
	MouseInput MouseInput
	UnitRepo   UnitRepo
}

func (sw *SpawnWarrior) Run() error {
	if !sw.MouseInput.IsPressed(input.MouseButtonLeft) {
		return nil
	}

	id := sw.IDGen.New()
	sw.UnitRepo.Upsert(id, model.Unit{
		Position: sw.MouseInput.Position(),
	})

	return nil
}
