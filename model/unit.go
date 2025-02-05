package model

import (
	"image"

	gmodel "github.com/elemir/gloomo/model"
)

type Unit struct {
	Animation *gmodel.Animation
	Position  image.Point
}
