package system

import (
	"github.com/elemir/gloomo/geom"
)

type Separation struct {
	UnitRepo UnitRepo
}

const (
	MinSeparation   = 16
	MaxSeparation   = 24
	SeparationCoeff = 1.5
)

// but later it should be a groups that moves in one direction.
func (s *Separation) Run() error {
	for id, unit := range s.UnitRepo.List() {
		var moveAway geom.Vec2
		var neighbours float64

		for otherID, otherUnit := range s.UnitRepo.List() {
			if id == otherID {
				continue
			}

			distance := unit.Position.Distance(otherUnit.Position)
			if distance < MaxSeparation {
				pushForce := unit.Position.Sub(otherUnit.Position)
				separationStrength := pushForce.Length() * (1 - SmoothStep(MinSeparation, MaxSeparation, distance))
				pushForce = pushForce.Normalize().Mul(pushForce.Length() * separationStrength)

				moveAway = moveAway.Add(pushForce)
				neighbours++
			}
		}

		if neighbours != 0 {
			moveAway = moveAway.Div(neighbours).Div(SeparationCoeff)
		}

		unit.Accel = unit.Accel.Add(moveAway)

		s.UnitRepo.Upsert(id, unit)
	}

	return nil
}

func SmoothStep(edge0, edge1, x float64) float64 {
	t := Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)

	return t * t * (3 - 2*t)
}

func Clamp(x, low, high float64) float64 {
	if x < low {
		return low
	}

	if x > high {
		return high
	}

	return x
}
