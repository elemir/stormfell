package main

import (
	"log/slog"
	"os"

	"github.com/elemir/gloomo"
	"github.com/elemir/gloomo/container"
	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/gloomo/loader"
	gmodel "github.com/elemir/gloomo/model"
	"github.com/elemir/gloomo/repo"
	"github.com/elemir/stormfell/algo"
	"github.com/elemir/stormfell/model"
	"github.com/elemir/stormfell/start"
	"github.com/hajimehoshi/ebiten/v2"
)

type Manager interface {
	Run() error
}

type Render interface {
	Draw(img *ebiten.Image)
}

type Game struct {
	w, h    int
	render  Render
	manager Manager
}

func (g *Game) Draw(img *ebiten.Image) {
	g.render.Draw(img)
}

func (g *Game) Update() error {
	err := g.manager.Run()
	if err != nil {
		slog.Error("Manager cycle run", slog.Any("err", err))
	}

	return nil
}

func (g *Game) Layout(w int, h int) (int, int) {
	g.w, g.h = w/3, h/3

	return g.w, g.h
}

func prepareManager(idGen *gid.Generator, tileMap *container.Resource[model.TileMap],
	spriteRepo *repo.Sprite, imgAssets *loader.Assets[*ebiten.Image],
) *gloomo.Manager {
	var manager gloomo.Manager

	perlinNoise := algo.NewPerlinNoise()
	fractalNoise := algo.NewFractalNoise(perlinNoise, 8, 0.5)

	manager.AddStartup(&start.MapGenerator{
		Coeff:  0.05,
		Width:  1000,
		Height: 1000,
		Noise:  fractalNoise,
		Levels: []float64{-0.3, 0},
		Setter: tileMap,
	})

	manager.AddStartup(&start.CreateTiles{
		IDGen:       idGen,
		SpriteRepo:  spriteRepo,
		Getter:      tileMap,
		ImageLoader: imgAssets,
	})

	manager.Add(&loader.Image{
		AssetDir: "assets",
		Assets:   imgAssets,
	})

	return &manager
}

func main() {
	var tileMap container.Resource[model.TileMap]
	var nodes container.SparseArray[gmodel.Node]
	var images container.SparseArray[*ebiten.Image]
	var imgAssets loader.Assets[*ebiten.Image]
	var idGen gid.Generator

	nodeRepo := &repo.Node{
		Nodes: &nodes,
	}

	spriteRepo := &repo.Sprite{
		Nodes:  &nodes,
		Images: &images,
	}

	rend := gloomo.NewRender(nodeRepo)
	manager := prepareManager(&idGen, &tileMap, spriteRepo, &imgAssets)

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(&Game{
		manager: manager,
		render:  rend,
	}); err != nil {
		slog.Error("Unable to run game", slog.Any("err", err))
		os.Exit(-1)
	}
}
