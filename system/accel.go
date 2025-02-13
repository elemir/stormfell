package system

import (
	"math"

	"github.com/elemir/gloomo/geom"
)

const (
	MaxSpeed  = 1
	MinSpeed  = 0.2
	MaxAngle  = math.Pi / 20
	SpeedDamp = 0.9
)

type Accel struct {
	UnitRepo UnitRepo
}

func (a *Accel) Run() error {
	for id, unit := range a.UnitRepo.List() {
		velocity := unit.Velocity.Add(unit.Accel)

		/*
			angle := unit.Velocity.AngleBetween(velocity)

			/*
			if angle > MaxAngle || angle < -MaxAngle {
				angle = geom.Angle(math.Copysign(MaxAngle, float64(angle)))
				velocity = unit.Velocity.Rotate(angle).Strip()
			}
			/*
				if angle < MinAngle && angle > -MinAngle {
					velocity = unit.Velocity.ResetLength(velocity.Length)
				}
		*/
		if velocity.Length() > MaxSpeed {
			velocity = velocity.Normalize().Mul(MaxSpeed)
		} else if velocity.Length() < MinSpeed {
			velocity = geom.Vec2{}
		}

		unit.Velocity = velocity.Mul(SpeedDamp)
		unit.Accel = geom.Vec2{}

		a.UnitRepo.Upsert(id, unit)
	}

	return nil
}
