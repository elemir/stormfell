package main

import (
	"log/slog"
	"os"

	"github.com/elemir/gloomo"
	"github.com/elemir/gloomo/container"
	"github.com/elemir/gloomo/geom"
	gid "github.com/elemir/gloomo/id"
	"github.com/elemir/gloomo/input"
	"github.com/elemir/gloomo/loader"
	gmodel "github.com/elemir/gloomo/model"
	"github.com/elemir/gloomo/node"
	grepo "github.com/elemir/gloomo/repo"
	gsystem "github.com/elemir/gloomo/system"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elemir/stormfell/algo"
	"github.com/elemir/stormfell/model"
	"github.com/elemir/stormfell/repo"
	"github.com/elemir/stormfell/start"
	"github.com/elemir/stormfell/system"
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
	g.w, g.h = 640, h*640/w

	return g.w, g.h
}

func prepareManager(spriteRepo *grepo.Sprite) *gloomo.Manager {
	var idGen gid.Generator
	var tileMap container.Resource[model.TileMap]

	var objPositions container.SparseArray[geom.Vec2]
	var objVelocities container.SparseArray[geom.Vec2]
	var objAccels container.SparseArray[geom.Vec2]
	var animations container.SparseArray[*gmodel.AnimationSheet]
	var stepCounters container.SparseArray[int]
	var currentAnimations container.SparseArray[string]
	var stoppedAnimations container.Set
	var zIndices container.SparseArray[int]

	var mouseInput input.Mouse
	var imgAssets loader.Assets[*ebiten.Image]
	var animAssets loader.Assets[*gmodel.AnimationSheet]
	var manager gloomo.Manager

	perlinNoise := algo.NewPerlinNoise()
	fractalNoise := algo.NewFractalNoise(perlinNoise, 8, 0.5)

	unitRepo := &repo.Unit{
		Animations:        &animations,
		Positions:         &objPositions,
		Velocities:        &objVelocities,
		Accelerations:     &objAccels,
		ZIndices:          &zIndices,
		CurrentAnimations: &currentAnimations,
		StoppedAnimations: &stoppedAnimations,
	}

	animRepo := &grepo.AnimatedSprite{
		Animations:        &animations,
		Positions:         &objPositions,
		ZIndices:          &zIndices,
		StepCounters:      &stepCounters,
		CurrentAnimations: &currentAnimations,
		StoppedAnimations: &stoppedAnimations,
	}

	manager.AddStartup(&start.MapGenerator{
		Coeff:  0.05,
		Width:  100,
		Height: 100,
		Noise:  fractalNoise,
		Levels: []float64{-0.3, 0},
		Setter: &tileMap,
	})

	manager.AddStartup(&start.CreateTiles{
		IDGen:       &idGen,
		SpriteRepo:  spriteRepo,
		Getter:      &tileMap,
		ImageLoader: &imgAssets,
	})

	manager.Add(&loader.Image{
		AssetDir: "assets",
		Assets:   &imgAssets,
	})

	manager.Add(&loader.Animation{
		AssetDir: "assets",
		Assets:   &animAssets,
	})

	manager.Add(&system.SpawnWarrior{
		IDGen:           &idGen,
		MouseInput:      &mouseInput,
		UnitRepo:        unitRepo,
		AnimationLoader: &animAssets,
	})

	manager.Add(&system.Separation{
		UnitRepo: unitRepo,
	})

	manager.Add(&system.Cohesion{
		UnitRepo: unitRepo,
	})

	manager.Add(&system.Alignment{
		UnitRepo: unitRepo,
	})

	manager.Add(&system.WallAvoid{
		TileMap:  &tileMap,
		UnitRepo: unitRepo,
	})

	manager.Add(&system.Accel{
		UnitRepo: unitRepo,
	})

	manager.Add(&system.Move{
		UnitRepo: unitRepo,
	})

	manager.Add(&gsystem.Animate{
		AnimationRepo: animRepo,
		SpriteRepo:    spriteRepo,
	})

	return &manager
}

func main() {
	var nodes container.SparseArray[node.Node]
	var images container.SparseArray[*ebiten.Image]
	var mirrors container.Set

	nodeRepo := &grepo.Node{
		Nodes: &nodes,
	}

	spriteRepo := &grepo.Sprite{
		Nodes:   &nodes,
		Images:  &images,
		Mirrors: &mirrors,
	}

	rend := gloomo.NewRender(nodeRepo)
	manager := prepareManager(spriteRepo)

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(&Game{
		manager: manager,
		render:  rend,
	}); err != nil {
		slog.Error("Unable to run game", slog.Any("err", err))
		os.Exit(-1)
	}
}
