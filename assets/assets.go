package assets

import (
	"embed"
	"image"
	_ "image/png" // Needed for PNG decoding, otherwise it panics
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *
var assets embed.FS

// NOTE: The above embed will show error if there is an empty subdirectory in 'assets'. Just make a dummy file to get rid of it

var FontFaceSource = mustLoadFontFace("fonts/antiquity-print.ttf")

var PlayerSprite = mustLoadSprite("tiles/1bit_tile_pack/Tilesheets/tilesheet_1of2_transparent.png", 16, 20, 17)
var EnemySprite = mustLoadSprite("tiles/1bit_tile_pack/Tilesheets/tilesheet_1of2_transparent.png", 16, 0, 0)

var (
	WallSprite_TopLeft = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile0.png",
	)
	WallSprite_Top = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile8.png",
	)
	WallSprite_TopRight = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile16.png",
	)
	WallSprite_Left = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile1.png",
	)
	WallSprite_Middle = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile9.png",
	)
	WallSprite_Right = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile17.png",
	)
	WallSprite_LowerLeft = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile2.png",
	)
	WallSprite_Lower = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile10.png",
	)
	WallSprite_LowerRight = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile18.png",
	)
)

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		log.Fatalf("Failed to open image: %v", name)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", name)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadSprite(name string, sz, row, col int) *ebiten.Image {
	x := col * sz
	y := row * sz
	return mustLoadImage(name).SubImage(image.Rect(x, y, x+sz, y+sz)).(*ebiten.Image)
}

func mustLoadFontFace(name string) *text.GoTextFaceSource {
	f, err := assets.Open(name)
	if err != nil {
		log.Fatalf("Failed to open font: %v", name)
	}
	defer f.Close()

	s, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatalf("Failed to load font: %v", name)
	}
	return s
}
