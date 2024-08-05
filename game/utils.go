package game

import (
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/Rabeez/rogue/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/teacat/noire"
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

func Abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func Sign(a int) int {
	if a >= 0 {
		return 1
	}
	return -1
}

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func NoireToColor(c noire.Color) color.Color {
	return color.RGBA{uint8(c.Red), uint8(c.Green), uint8(c.Blue), uint8(c.Alpha)}
}
