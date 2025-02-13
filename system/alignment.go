package system

import (
	"github.com/elemir/gloomo/geom"
)

type Alignment struct {
	UnitRepo UnitRepo
}

const (
	AlignmentCoeff  = 1000
	AlignmentRadius = 100
)

func (a *Alignment) Run() error {
	for id, unit := range a.UnitRepo.List() {
		var force geom.Vec2
		var neighbours float64

		for otherID, otherUnit := range a.UnitRepo.List() {
			if id == otherID {
				continue
			}

			distance := unit.Position.Distance(otherUnit.Position)
			if distance < AlignmentRadius {
				force = force.Add(unit.Velocity)
				neighbours++
			}
		}

		if neighbours == 0 {
			continue
		}

		force = force.Div(neighbours).Div(AlignmentCoeff)

		unit.Accel = unit.Accel.Add(force)

		a.UnitRepo.Upsert(id, unit)
	}

	return nil
}
