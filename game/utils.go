package game

import (
	"image/color"
	"math/rand/v2"

	"github.com/Rabeez/rogue/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/teacat/noire"
)

func PutText(panel *Panel, s string, op *text.DrawOptions, clr color.Color, sz float64) {
	f := &text.GoTextFace{
		Source: assets.FontFaceSource,
		Size:   sz,
	}
	// TODO: redo this using values from `op` instead
	// // Detect bounds and overflow drawing
	// w, h := text.Measure(s, f, 0)
	// if x+w > float64(panel.Size.X) || y+h > float64(panel.Size.Y) {
	// 	log.Fatalf("Overflow drawing detected\n")
	// }

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
