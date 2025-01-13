package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/elemir/gloomo"
	"github.com/elemir/gloomo/container"
	"github.com/elemir/gloomo/draw"
	gid "github.com/elemir/gloomo/id"
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
		slog.Error("unable to run manager", slog.Any("err", err))

		return fmt.Errorf("run manager: %w", err)
	}

	return nil
}

func (g *Game) Layout(w int, h int) (int, int) {
	g.w, g.h = w, h

	return w, h
}

func prepareManager(idGen *gid.Generator, tileMap *container.Resource[model.TileMap],
	nodeRepo *repo.Node,
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
		IDGen:    idGen,
		NodeRepo: nodeRepo,
		Getter:   tileMap,
	})

	return &manager
}

func main() {
	var tileMap container.Resource[model.TileMap]

	var nodes container.SparseArray[draw.Node]

	var idGen gid.Generator

	nodeRepo := &repo.Node{
		Nodes: &nodes,
	}

	rend := gloomo.NewRender(nodeRepo)
	manager := prepareManager(&idGen, &tileMap, nodeRepo)

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(&Game{
		manager: manager,
		render:  rend,
	}); err != nil {
		slog.Error("Unable to run game", slog.Any("err", err))
		os.Exit(-1)
	}
}
