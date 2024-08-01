package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/Rabeez/rogue/assets"
)

const (
	TITLE         = "Rogue"
	TILE_SIZE     = 16
	WINDOW_FACTOR = 4
	WINDOW_WIDTH  = 16 * WINDOW_FACTOR * TILE_SIZE
	WINDOW_HEIGHT = 9 * WINDOW_FACTOR * TILE_SIZE
)

type Game struct {
	player  *Player
	enemies []*Enemy
	walls   []*Wall
}

func NewGame() *Game {

	// wall layout from map file?
	w := []*Wall{NewWall()}

	// player initial location

	// enemy initial locations

	// interactables locations

	g := &Game{
		// player:  NewPlayer(),
		// enemies: []*Enemy{},
		walls: w,
	}
	return g
}

func (g *Game) Update() error {

	// keypresses + state update

	// enemy AI

	// use updated state from now:

	// interactions (chest open, pickups, open walls/doors)

	// attacks (projectiles?)

	// collisions? (here or with state updaes?)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")

	// draw background
	screen.Fill(color.RGBA{0x20, 0x20, 0x20, 0xff})

	// draw walls
	for _, w := range g.walls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(w.X, w.Y)
		screen.DrawImage(assets.WallTSprite, op)
	}

	// draw interactables (chests, items)

	// draw enemies

	// draw player

	// draw UI
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
