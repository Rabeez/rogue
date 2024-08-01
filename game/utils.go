package game

import (
	"image/color"

	"github.com/Rabeez/rogue/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawText(screen *ebiten.Image, s string, x, y float64, clr color.Color, sz float64) {
	f := &text.GoTextFace{
		Source: assets.FontFaceSource,
		Size:   sz,
	}
	// w, h := text.Measure(s, f, 0)
	// vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{200, 0, 0, 255}, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, s, f, op)
}
