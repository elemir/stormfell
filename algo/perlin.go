package algo

import (
	"math"
	"math/rand/v2"
)

const (
	permCount = 256
)

type PerlinNoise struct {
	perm []int
}

func NewPerlinNoise() *PerlinNoise {
	perm := make([]int, permCount*2)

	for i := range permCount {
		perm[i] = rand.IntN(permCount)
		perm[i+permCount] = perm[i]
	}

	return &PerlinNoise{
		perm: perm,
	}
}

func (pn *PerlinNoise) Noise(x, y float64) float64 {
	gridX, gridY := gridCoord(x), gridCoord(y)

	localX, localY := localCoord(x), localCoord(y)

	fadedX, fadedY := fade(localX), fade(localY)

	// Hash of the grid cell vertexes
	aa := pn.hash(gridX, gridY)
	ab := pn.hash(gridX, gridY+1)
	ba := pn.hash(gridX+1, gridY)
	bb := pn.hash(gridX+1, gridY+1)

	lowerX := lerp(grad(aa, localX, localY), grad(ba, localX-1, localY), fadedX)
	upperX := lerp(grad(ab, localX, localY-1), grad(bb, localX-1, localY-1), fadedX)

	return lerp(lowerX, upperX, fadedY)
}

func (pn *PerlinNoise) hash(gridX, gridY int) int {
	return pn.perm[pn.perm[gridX]+gridY]
}

func gridCoord(x float64) int {
	return int(math.Floor(x)) & (permCount - 1)
}

func localCoord(x float64) float64 {
	return x - math.Floor(x)
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	switch hash & 3 {
	case 0:
		return x + y
	case 1:
		return y - x
	case 2:
		return x - y
	case 3:
		return -x - y
	}

	return 0
}
