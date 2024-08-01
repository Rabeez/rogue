package assets

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var assets embed.FS

var PlayerSprite = mustLoadImage("tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile1.png")
var EnemySprite = mustLoadImage("tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile1.png")
var WallTSprite = mustLoadImage("tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile1.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
