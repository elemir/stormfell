package start

import (
	"errors"
	"image"

	"github.com/elemir/gloomo/draw"
	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/stormfell/model"
)

var errTileMapIsNotExist = errors.New("tile map is not exist")

type IDGenerator interface {
	New() gid.ID
}

type TileMapGetter interface {
	Get() (model.TileMap, bool)
}

type NodeRepo interface {
	Get(id gid.ID) (draw.Node, bool)
	Upsert(id gid.ID, node draw.Node)
}

type CreateTiles struct {
	IDGen    IDGenerator
	Getter   TileMapGetter
	NodeRepo NodeRepo
}

func (ct *CreateTiles) Run() error {
	tileMap, ok := ct.Getter.Get()
	if !ok {
		return errTileMapIsNotExist
	}

	drawRect := draw.Rect(ct.NodeRepo)

	for i := range 50 {
		for j := range 50 {
			if tileMap[i][j] == 1 {
				id := ct.IDGen.New()

				ct.NodeRepo.Upsert(id, draw.Node{
					Draw:     drawRect,
					Position: image.Pt(i*64, j*64),
					Size:     image.Pt(64, 64),
				})
			}
		}
	}

	return nil
}
