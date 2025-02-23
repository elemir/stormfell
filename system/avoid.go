package system

import (
	"errors"
	"image"

	"github.com/elemir/gloomo/geom"

	"github.com/elemir/stormfell/model"
)

const (
	AvoidCoeff = 200
)

var errTileMapIsNotExist = errors.New("tile map is not exist")

type TileMapGetter interface {
	Get() (model.TileMap, bool)
}

type WallAvoid struct {
	TileMap  TileMapGetter
	UnitRepo UnitRepo
}

func (wa *WallAvoid) Run() error {
	tileMap, ok := wa.TileMap.Get()
	if !ok {
		return errTileMapIsNotExist
	}

	for id, unit := range wa.UnitRepo.List() {
		var force geom.Vec2
		pt := unit.Position.Div(32).Round()

		for neighPt, val := range tileMap.Neighbours(pt.X, pt.Y) {
			if val != 1 {
				continue
			}

			// TODO(evgenii.omelchenko): we should have a service that knows how to convert tiles into real position
			neighCenter := geom.FromPoint(neighPt.Mul(32).Add(image.Pt(16, 16)))

			pushForce := unit.Position.Sub(neighCenter)
			distance := pushForce.Length()

			if distance < MaxSeparation {
				separationStrength := pushForce.Length() * (1 - SmoothStep(MinSeparation, MaxSeparation, distance))
				pushForce = pushForce.Normalize().Mul(pushForce.Length() * separationStrength)

				force = force.Add(pushForce)
			}
		}

		unit.Accel = unit.Accel.Add(force).Div(AvoidCoeff)

		wa.UnitRepo.Upsert(id, unit)
	}

	return nil
}
