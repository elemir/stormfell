package start

import (
	"errors"
	"image"
	"math/rand/v2"

	gid "github.com/elemir/gloomo/id"
	gmodel "github.com/elemir/gloomo/model"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/stormfell/model"
)

var (
	errTileMapIsNotExist = errors.New("tile map is not exist")
	errImageNotLoaded    = errors.New("image not loaded")
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
	Get(id gid.ID) (gmodel.Sprite, bool)
	Upsert(id gid.ID, sprite gmodel.Sprite)
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

	waterImg, ok := ct.ImageLoader.Load("water.png")
	if !ok {
		return errImageNotLoaded
	}

	for i := range 50 {
		for j := range 50 {
			img := randImg(waterImg, 5)

			if tileMap[i][j] == 1 {
				img = randImg(dirtImg, 3)
			}

			id := ct.IDGen.New()
			ct.SpriteRepo.Upsert(id, gmodel.Sprite{
				Image:    img,
				Position: image.Pt(i*32, j*32),
			})
		}
	}

	return nil
}

func randImg(img *ebiten.Image, count int) *ebiten.Image {
	n := rand.IntN(count)

	//nolint:forcetypeassert // SubImage for ebiten image always returns ebiten image
	return img.SubImage(image.Rect(n*32, 0, n*32+32, 32)).(*ebiten.Image)
}
