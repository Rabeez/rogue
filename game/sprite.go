package game

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/Rabeez/rogue/assets"
)

func make_speed(speed float64) float64 {
	// Convert speed in tiles-per-second to ticks-per-second
	return speed * float64(TILE_SIZE) / float64(ebiten.TPS())
}

type Sprite struct {
	// Sprite coords are in grid coord space
	X, Y  float64
	Img   *ebiten.Image
	color color.Color
}

func (p *Sprite) Draw(screen *ebiten.Image) {
	// This method converts sprite grid coords to screen pixel coords for drawing
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)
	op.ColorScale.ScaleWithColor(p.color)
	screen.DrawImage(p.Img, op)
}

type Player struct {
	*Sprite
	speed float64
}

func NewPlayer(x, y float64) *Player {
	// Move 300 pixels per second
	s := make_speed(30)
	return &Player{
		speed: s,
		Sprite: &Sprite{
			X:     x,
			Y:     y,
			color: color.RGBA{0xff, 0xff, 0x00, 0xff},
			Img:   assets.PlayerSprite,
		},
	}
}

func (p *Player) Update() {
	var deltaX, deltaY float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		deltaY = -p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		deltaY = p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		deltaX = -p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		deltaX = p.speed
	}

	// Check for diagonal movement
	if deltaX != 0 && deltaY != 0 {
		factor := p.speed / math.Sqrt(deltaX*deltaX+deltaY*deltaY)
		deltaX *= factor
		deltaY *= factor
	}

	p.X += deltaX
	p.Y += deltaY
}

type Enemy struct {
	*Sprite
}

func NewEnemy(x, y float64) *Enemy {
	return &Enemy{
		Sprite: &Sprite{
			X:     x,
			Y:     y,
			color: color.RGBA{0xaa, 0x20, 0x20, 0xff},
			Img:   assets.EnemySprite,
		},
	}
}

type WallType int

const (
	Wall_TopLeft WallType = iota
	Wall_Top
	Wall_TopRight
	Wall_Left
	Wall_Middle
	Wall_Right
	Wall_LowerLeft
	Wall_Lower
	Wall_LowerRight
)

type Wall struct {
	*Sprite
	wallType WallType
}

func NewWall(x, y float64, wallType WallType) *Wall {
	sp := &Sprite{
		X:     x,
		Y:     y,
		color: color.White,
		Img:   assets.WallSprite_TopLeft,
	}
	switch wallType {
	case Wall_TopLeft:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Top:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_TopRight:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Left:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Middle:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Right:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LowerLeft:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Lower:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LowerRight:
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	default:
		log.Fatalf("Invalid wall type provided: %d", wallType)
		return nil
	}
}
