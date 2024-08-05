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

// TODO: should TILE_SIZE be used here, or should these be manually resized to match?
var PlayerSprite = mustLoadImage("tiles/1bit_tile_pack/extracted/1_01-01.png")
var EnemySprite = mustLoadImage("tiles/1bit_tile_pack/extracted/1_02-01.png")
var CoinSprite = mustLoadImage("tiles/1bit_tile_pack/extracted/1_00-11.png")
var ChestSprite = mustLoadImage("tiles/1bit_tile_pack/extracted/1_06-10.png")

var (
	WallSprite_TopLeft = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile0.png",
	)
	WallSprite_TopT = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile8.png",
	)
	WallSprite_TopRight = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile16.png",
	)
	WallSprite_LeftT = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile1.png",
	)
	WallSprite_Middle = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile9.png",
	)
	WallSprite_RightT = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile17.png",
	)
	WallSprite_LowerLeft = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile2.png",
	)
	WallSprite_LowerT = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile10.png",
	)
	WallSprite_LowerRight = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile18.png",
	)
	WallSprite_Horz = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile57.png",
	)
	WallSprite_Vert = mustLoadImage(
		"tiles/Ascii-wall-like tileset by GibberishDev/connected/connected-tile56.png",
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
