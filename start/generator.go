package start

import (
	"maps"

	"github.com/elemir/stormfell/algo"
	"github.com/elemir/stormfell/model"
)

type Noise interface {
	Noise(x, y float64) float64
}

type TileMapSetter interface {
	Set(tm model.TileMap)
}

type MapGenerator struct {
	Coeff         float64
	Width, Height int
	Levels        []float64
	Noise         Noise
	Setter        TileMapSetter
}

func (mg *MapGenerator) Run() error {
	tileMap := make(model.TileMap, mg.Width)

	for i := range mg.Width {
		tileMap[i] = make([]int, mg.Height)

		for j := range mg.Height {
			noise := mg.Noise.Noise(float64(i)*mg.Coeff, float64(j)*mg.Coeff)

			tileMap[i][j] = algo.BinSearch(mg.Levels, noise)
		}
	}

	for range 5 {
		for i := range mg.Width {
			for j := range mg.Height {
				if len(maps.Collect(tileMap.Neighbours(i, j))) <= 3 {
					tileMap[i][j] = 0
				}
			}
		}
	}

	mg.Setter.Set(tileMap)

	return nil
}
