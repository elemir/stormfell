package system

import (
	"slices"

	"github.com/elemir/gloomo/geom"

	"github.com/elemir/stormfell/model"
)

const (
	PerceptionRadius = 200
	CohesionCoeff    = 250
	MaxNeighbours    = 5
)

type Cohesion struct {
	UnitRepo UnitRepo
}

func (c *Cohesion) Run() error {
	for id, unit := range c.UnitRepo.List() {
		var centerOfMass geom.Vec2
		var all []model.Unit
		var neighbours float64

		// TODO(evgenii.omelchenko): use heap here
		for otherID, otherUnit := range c.UnitRepo.List() {
			if id == otherID {
				continue
			}

			all = append(all, otherUnit)
		}

		slices.SortFunc(all, func(u, v model.Unit) int {
			return int(unit.Position.Distance(u.Position) - unit.Position.Distance(v.Position))
		})

		for _, otherUnit := range all {
			distance := unit.Position.Distance(otherUnit.Position)
			if distance < PerceptionRadius && distance > MaxSeparation {
				centerOfMass = centerOfMass.Add(otherUnit.Position)
				neighbours++
			}

			if neighbours >= MaxNeighbours {
				break
			}
		}

		var force geom.Vec2

		if neighbours != 0 {
			centerOfMass = centerOfMass.Div(neighbours)

			force = centerOfMass.Sub(unit.Position).Div(CohesionCoeff)
		}

		unit.Accel = unit.Accel.Add(force)

		c.UnitRepo.Upsert(id, unit)
	}

	return nil
}
