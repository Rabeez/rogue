package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Panel struct {
	Corner Vector2
	Size   Vector2
	Screen *ebiten.Image
}

func NewPanel(x, y, w, h int, img *ebiten.Image) *Panel {
	// All 4 coords are in absolute screen space coords
	return &Panel{
		Corner: *NewVector2(x, y),
		Size:   *NewVector2(w, h),
		Screen: img.SubImage(image.Rect(x, y, w, h)).(*ebiten.Image),
	}
}

func (p *Panel) SubPanel(rect image.Rectangle) *Panel {
	// Rect coords are relative to the Panel p (not absolute screen coords)
	// So this method converts rect coords to absolute screen coords
	return NewPanel(p.Corner.X+rect.Min.X, p.Corner.Y+rect.Min.Y, p.Size.X+rect.Dx(), p.Size.Y+rect.Dy(), p.Screen)
}
