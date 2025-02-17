package model

import (
	"image"
	"iter"
)

type TileMap [][]int

func (t TileMap) Size() (int, int) {
	if len(t) == 0 {
		return 0, 0
	}

	return len(t), len(t[0])
}

func (t TileMap) Neighbours(x0, y0 int) iter.Seq2[image.Point, int] {
	return func(yield func(image.Point, int) bool) {
		for i := range 3 {
			x := x0 + i - 1

			for j := range 3 {
				y := y0 + j - 1

				isCenter := x == x0 && y == y0
				val, insideBorder := t.At(x, y)

				if isCenter || !insideBorder {
					continue
				}

				if !yield(image.Pt(x, y), val) {
					return
				}
			}
		}
	}
}

func (t TileMap) At(x, y int) (int, bool) {
	width, height := t.Size()

	outOfBorder := x < 0 || y < 0 || x >= width || y >= height
	if outOfBorder {
		return 0, false
	}

	return t[x][y], true
}
