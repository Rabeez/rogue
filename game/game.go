package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TITLE     = "Rogue"
	TILE_SIZE = 16
	// WINDOW_FACTOR = 4
	WINDOW_WIDTH  = 80 * TILE_SIZE
	WINDOW_HEIGHT = 50 * TILE_SIZE
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
		// TODO: pause game on escape and show menu with exit option instead
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
	// Setup canvases
	gamePanel := NewPanel(0, 0, WINDOW_HEIGHT, WINDOW_HEIGHT, screen)
	interfacePanel := NewPanel(WINDOW_HEIGHT, 0, WINDOW_WIDTH, WINDOW_HEIGHT, screen)
	interfaceLeftPanel := interfacePanel.SubPanel(image.Rect(0, 0, 100, WINDOW_HEIGHT))
	interfaceRightPanel := interfacePanel.SubPanel(image.Rect(100, 0, WINDOW_WIDTH, WINDOW_HEIGHT))

	// Draw game level w/ entitites
	g.currentLevel.Draw(gamePanel)

	// Draw UI
	PutText(interfaceLeftPanel, "hello", 0, 0, color.RGBA{200, 0, 0, 255}, 24)
	PutText(interfaceRightPanel, "hello again", 100, 100, color.RGBA{10, 0, 100, 255}, 30)
	PutText(interfaceRightPanel, fmt.Sprintf("Gold: %d", g.currentLevel.Player.gold), 100, 130, color.RGBA{10, 50, 100, 255}, 30)

	// Debug info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Return hardcoded values and allow engine to redraw appropriatel
	// https://ebitengine.org/en/blog/resizable.html
	return WINDOW_WIDTH, WINDOW_HEIGHT
}
