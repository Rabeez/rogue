package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TITLE         = "Rogue"
	TILE_SIZE     = 16
	WINDOW_WIDTH  = 80 * TILE_SIZE
	WINDOW_HEIGHT = 50 * TILE_SIZE
)

var (
	BACKGROUND_COLOR = color.RGBA{0x20, 0x25, 0x20, 0xff}
)

type Game struct {
	isPaused     bool
	currentLevel *Level
}

func NewGame() *Game {
	g := &Game{
		isPaused:     false,
		currentLevel: NewLevel(1),
	}
	return g
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.isPaused = !g.isPaused
	}

	if !g.isPaused {
		g.currentLevel.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Setup canvases
	gamePanel := NewPanel(0, 0, WINDOW_HEIGHT, WINDOW_HEIGHT, screen)
	_interfacePanel := NewPanel(WINDOW_HEIGHT, 0, WINDOW_WIDTH, WINDOW_HEIGHT, screen)
	interfaceLeftPanel := _interfacePanel.SubPanel(image.Rect(0, 0, 100, WINDOW_HEIGHT))
	interfaceRightPanel := _interfacePanel.SubPanel(image.Rect(100, 0, WINDOW_WIDTH, WINDOW_HEIGHT))

	// Draw game level w/ entitites
	g.currentLevel.Draw(gamePanel)

	// Draw UI
	PutText(interfaceLeftPanel, "hello", 0, 0, color.RGBA{200, 0, 0, 255}, 24)
	PutText(interfaceRightPanel, "hello again", 100, 100, color.RGBA{10, 0, 100, 255}, 30)
	PutText(interfaceRightPanel, fmt.Sprintf("Gold: %d", g.currentLevel.Player.gold), 100, 130, color.RGBA{10, 50, 100, 255}, 30)
	PutText(interfaceRightPanel, fmt.Sprintf("Health: %d", g.currentLevel.Player.health), 100, 200, color.RGBA{100, 50, 100, 255}, 30)

	// Debug info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Return hardcoded values and allow engine to redraw appropriatel
	// https://ebitengine.org/en/blog/resizable.html
	return WINDOW_WIDTH, WINDOW_HEIGHT
}
