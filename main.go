package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/Rabeez/rogue/game"
)

func main() {
	ebiten.SetWindowSize(game.WINDOW_WIDTH, game.WINDOW_HEIGHT)
	ebiten.SetWindowTitle(game.TITLE)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
