package model

import (
	"github.com/elemir/gloomo/geom"
	gmodel "github.com/elemir/gloomo/model"
)

type Unit struct {
	Animation *gmodel.AnimationSheet
	Position  geom.Vec2
	Velocity  geom.Vec2
}
