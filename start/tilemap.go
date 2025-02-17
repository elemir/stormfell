package start

import (
	"errors"
	"image"
	"math/rand/v2"

	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/gloomo/node"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/stormfell/model"
)

var (
	errTileMapIsNotExist = errors.New("tile map is not exist")
	errImageNotLoaded    = errors.New("image not loaded")

	// TODO(evgenii.omelchenko): should be loaded from yaml.
	rock = model.Terrain{
		Transitions: map[string][]int{
			"solid": {18, 21, 24},

			"north": {2, 5, 8},
			"south": {34, 37, 40},
			"west":  {17, 20, 23},
			"east":  {19, 22, 25},

			"northwest_outer": {1, 4},
			"northeast_outer": {3, 6},
			"southwest_outer": {33, 36},
			"southeast_outer": {35, 38},

			"northwest_inner": {10, 12},
			"northeast_inner": {11, 13},
			"southwest_inner": {26, 28},
			"southeast_inner": {27, 29},

			"northwest_southeast_inner": {30, 31},
			"northeast_southwest_inner": {14, 15},
		},
	}
)

type IDGenerator interface {
	New() gid.ID
}

type TileMapGetter interface {
	Get() (model.TileMap, bool)
}

type ImageLoader interface {
	Load(path string) (*ebiten.Image, bool)
}

type SpriteRepo interface {
	Get(id gid.ID) (node.Sprite, bool)
	Upsert(id gid.ID, sprite node.Sprite)
}

type CreateTiles struct {
	IDGen       IDGenerator
	Getter      TileMapGetter
	SpriteRepo  SpriteRepo
	ImageLoader ImageLoader
}

func (ct *CreateTiles) Run() error {
	tileMap, ok := ct.Getter.Get()
	if !ok {
		return errTileMapIsNotExist
	}

	dirtImg, ok := ct.ImageLoader.Load("dirt.png")
	if !ok {
		return errImageNotLoaded
	}

	rockImg, ok := ct.ImageLoader.Load("rock.png")
	if !ok {
		return errImageNotLoaded
	}

	for i := range 50 {
		for j := range 50 {
			img := randImg(dirtImg, 3)

			id := ct.IDGen.New()
			ct.SpriteRepo.Upsert(id, node.Sprite{
				Image:    img,
				Position: image.Pt(i*32, j*32),
			})

			if tileMap[i][j] == 1 {
				mask := tileMapTransitionMask(tileMap, i, j)
				img := randRock(rockImg, mask)

				id := ct.IDGen.New()
				ct.SpriteRepo.Upsert(id, node.Sprite{
					Image:    img,
					Position: image.Pt(i*32, j*32),
				})
			}
		}
	}

	return nil
}

func randImg(img *ebiten.Image, count int) *ebiten.Image {
	n := rand.IntN(count)

	//nolint:forcetypeassert // SubImage for ebiten image always returns ebiten image
	return img.SubImage(image.Rect(n*32, 0, n*32+32, 32)).(*ebiten.Image)
}

//nolint:cyclop // this function should be properly rewrited with using some abstract rules
func randRock(img *ebiten.Image, mask transitionMask) *ebiten.Image {
	transition := "solid"

	switch {
	case mask&0xB == 0xB:
		transition = "northwest_outer"
	case mask&0x16 == 0x16:
		transition = "southwest_outer"
	case mask&0x68 == 0x68:
		transition = "northeast_outer"
	case mask&0xD0 == 0xD0:
		transition = "southeast_outer"

	case mask == 0x1:
		transition = "northwest_inner"
	case mask == 0x4:
		transition = "southwest_inner"
	case mask == 0x20:
		transition = "northeast_inner"
	case mask == 0x80:
		transition = "southeast_inner"

	case mask == 0x81:
		transition = "northwest_southeast_inner"
	case mask == 0x24:
		transition = "northeast_southwest_inner"

	case mask&0x8 != 0:
		transition = "north"
	case mask&0x10 != 0:
		transition = "south"
	case mask&0x2 != 0:
		transition = "west"
	case mask&0x40 != 0:
		transition = "east"
	}

	n := rand.IntN(len(rock.Transitions[transition]))

	return extractSpriteImg(img, rock.Transitions[transition][n])
}

func extractSpriteImg(img *ebiten.Image, n int) *ebiten.Image {
	width := img.Bounds().Dx()
	countInRow := width / 32
	x, y := n%countInRow, n/countInRow

	rect := image.Rect(x*32, y*32, (x+1)*32, (y+1)*32)

	//nolint:forcetypeassert // SubImage for ebiten image always returns ebiten image
	return img.SubImage(rect).(*ebiten.Image)
}

// TODO(evgenii.omelchenko): looks rotated some reason.
var shift = map[image.Point]int{
	image.Pt(-1, -1): 0, // north-west (0x01)
	image.Pt(-1, 0):  1, // west (0x02)
	image.Pt(-1, 1):  2, // south-west (0x04)
	image.Pt(0, -1):  3, // north (0x08)
	image.Pt(0, 1):   4, // south (0x10)
	image.Pt(1, -1):  5, // north-east (0x20)
	image.Pt(1, 0):   6, // east (0x40)
	image.Pt(1, 1):   7, // south-east (0x80)
}

type transitionMask byte

func tileMapTransitionMask(tileMap model.TileMap, x, y int) transitionMask {
	var mask transitionMask

	pt := image.Pt(x, y)

	for neighbour, val := range tileMap.Neighbours(x, y) {
		if tileMap[x][y] != val {
			mask |= 0x1 << shift[neighbour.Sub(pt)]
		}
	}

	return mask
}
