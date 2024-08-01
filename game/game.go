package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TITLE         = "Rogue"
	TILE_SIZE     = 16
	WINDOW_FACTOR = 4
	WINDOW_WIDTH  = 16 * WINDOW_FACTOR * TILE_SIZE
	WINDOW_HEIGHT = 9 * WINDOW_FACTOR * TILE_SIZE
)

var (
	BACKGROUND_COLOR = color.RGBA{0x20, 0x25, 0x20, 0xff}
)

type Game struct {
	currentLevel *Level
}

func NewGame() *Game {

	// wall layout from map file?

	// player initial location
	// enemy initial locations
	// interactables locations

	g := &Game{
		currentLevel: NewLevel(1),
	}
	return g
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	g.currentLevel.Update()
	// keypresses + state update

	// enemy AI

	// use updated state from now:

	// interactions (chest open, pickups, open walls/doors)

	// attacks (projectiles?)

	// collisions? (here or with state updaes?)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentLevel.Draw(screen)
	// draw interactables (chests, items)

	// draw enemies

	// draw player

	// draw UI
	DrawText(screen, "hello", 300, 300, color.RGBA{200, 0, 0, 255}, 24)
	DrawText(screen, "hello again", 500, 500, color.RGBA{10, 0, 100, 255}, 30)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
