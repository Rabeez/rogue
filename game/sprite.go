package game

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/Rabeez/rogue/assets"
)

type Sprite struct {
	// Sprite coords are in grid coord space
	Pos   *Vector2
	Img   *ebiten.Image
	color color.Color
}

func (p *Sprite) Draw(screen *ebiten.Image) {
	// This method converts sprite grid coords to screen pixel coords for drawing
	op := &ebiten.DrawImageOptions{}
	pixelX := float64(p.Pos.X * TILE_SIZE)
	pixelY := float64(p.Pos.Y * TILE_SIZE)
	op.GeoM.Translate(pixelX, pixelY)
	op.GeoM.Scale(2, 2)
	op.ColorScale.ScaleWithColor(p.color)
	screen.DrawImage(p.Img, op)
}

type Player struct {
	*Sprite
	speed int
}

func NewPlayer(x, y int) *Player {
	// speed is in units of cells
	s := 1
	if s < 1 {
		log.Fatalf("Player speed has to be >= 1: %v", s)
	}
	return &Player{
		speed: s,
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xff, 0xff, 0x00, 0xff},
			Img:   assets.PlayerSprite,
		},
	}
}

func (p *Player) Update(colliders *map[Vector2]bool) {
	deltaPos := NewVector2(0, 0)
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		deltaPos.Y = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		deltaPos.Y = p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		deltaPos.X = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		deltaPos.X = p.speed
	}

	newPos := p.Pos.Add(*deltaPos)
	if _, ok := (*colliders)[newPos]; !ok {
		p.Pos = &newPos
	}
}

type Enemy struct {
	*Sprite
}

func NewEnemy(x, y int) *Enemy {
	return &Enemy{
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xaa, 0x20, 0x20, 0xff},
			Img:   assets.EnemySprite,
		},
	}
}

type WallType int

const (
	Wall_TopLeft WallType = iota
	Wall_TopT
	Wall_TopRight
	Wall_LeftT
	Wall_Middle
	Wall_RightT
	Wall_LowerLeft
	Wall_LowerT
	Wall_LowerRight
	Wall_Horz
	Wall_Vert
)

type Wall struct {
	*Sprite
	wallType WallType
}

func NewWall(x, y int, wallType WallType) *Wall {
	sp := &Sprite{
		Pos:   NewVector2(x, y),
		color: color.White,
		Img:   nil,
	}
	switch wallType {
	case Wall_TopLeft:
		sp.Img = assets.WallSprite_TopLeft
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_TopT:
		sp.Img = assets.WallSprite_TopT
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_TopRight:
		sp.Img = assets.WallSprite_TopRight
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LeftT:
		sp.Img = assets.WallSprite_LeftT
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Middle:
		sp.Img = assets.WallSprite_Middle
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_RightT:
		sp.Img = assets.WallSprite_RightT
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LowerLeft:
		sp.Img = assets.WallSprite_LowerLeft
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LowerT:
		sp.Img = assets.WallSprite_LowerT
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_LowerRight:
		sp.Img = assets.WallSprite_LowerRight
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Horz:
		sp.Img = assets.WallSprite_Horz
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	case Wall_Vert:
		sp.Img = assets.WallSprite_Vert
		return &Wall{
			wallType: wallType,
			Sprite:   sp,
		}
	default:
		log.Fatalf("Invalid wall type provided: %d", wallType)
		return nil
	}
}
