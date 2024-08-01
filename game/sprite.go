package game

import (
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
	X, Y float64
	Img  *ebiten.Image
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
			X:   x,
			Y:   y,
			Img: assets.PlayerSprite,
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

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.Img, op)
}

type Enemy struct {
	*Sprite
}

func NewEnemy(x, y float64) *Enemy {
	return &Enemy{
		Sprite: &Sprite{
			X:   x,
			Y:   y,
			Img: assets.EnemySprite,
		},
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(e.Img, op)
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
	switch wallType {
	case Wall_TopLeft:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_TopLeft,
			},
		}
	case Wall_Top:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Top,
			},
		}
	case Wall_TopRight:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_TopRight,
			},
		}
	case Wall_Left:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Left,
			},
		}
	case Wall_Middle:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Middle,
			},
		}
	case Wall_Right:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Right,
			},
		}
	case Wall_LowerLeft:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_LowerLeft,
			},
		}
	case Wall_Lower:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Lower,
			},
		}
	case Wall_LowerRight:
		return &Wall{
			wallType: wallType,
			Sprite: &Sprite{
				X:   x,
				Y:   y,
				Img: assets.WallSprite_Right,
			},
		}
	default:
		log.Fatalf("Invalid wall type provided: %d", wallType)
		return nil
	}
}

func (w *Wall) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.X, w.Y)
	screen.DrawImage(w.Img, op)
}
