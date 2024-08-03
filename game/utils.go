package game

import (
	"image/color"
	"log"

	"github.com/Rabeez/rogue/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func PutText(panel *Panel, s string, x, y float64, clr color.Color, sz float64) {
	f := &text.GoTextFace{
		Source: assets.FontFaceSource,
		Size:   sz,
	}
	// Detect bounds and overflow drawing
	w, h := text.Measure(s, f, 0)
	if x+w > float64(panel.Size.X) || y+h > float64(panel.Size.Y) {
		log.Fatalf("Overflow drawing detected\n")
	}

	op := &text.DrawOptions{}
	// Account for panel offset
	op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(panel.Screen, s, f, op)
}

// TODO: PutRect(filled bool)

// TODO: PutParticleEmitter(n int)

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func sign(a int) int {
	if a >= 0 {
		return 1
	}
	return -1
}
