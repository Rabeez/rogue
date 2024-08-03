package main

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/Rabeez/rogue/game"
)

func main() {
	ebiten.SetWindowSize(game.WINDOW_WIDTH, game.WINDOW_HEIGHT)
	ebiten.SetWindowTitle(game.TITLE)
	ebiten.SetVsyncEnabled(true)
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		if !errors.Is(err, ebiten.Termination) {
			log.Fatal(err)
		}
	}
}
