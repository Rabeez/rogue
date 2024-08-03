package game

import (
	"image/color"
	"log"
	"math"

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
	speed float64
}

func NewPlayer(x, y int) *Player {
	// speed is in units of cells
	s := 1.0
	if s < 1.0 {
		log.Fatalf("Player speed has to be >= 1.0: %v", s)
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

func (p *Player) Update(walls []*Wall) {
	// TODO: speed is too fast. maybe just fix it to grid cells and remove diagonal option
	var deltaX, deltaY float64

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		deltaY = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		deltaY = p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		deltaX = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		deltaX = p.speed
	}

	// Check for diagonal movement
	// if deltaX != 0 && deltaY != 0 {
	// 	factor := p.speed / math.Sqrt(deltaX*deltaX+deltaY*deltaY)
	// 	deltaX *= factor
	// 	deltaY *= factor
	// }

	// Check for wall collisions
	// possibleCollisionCoords := [][]int{}
	// for row_offset := -1; row_offset <= 1; row_offset++ {
	// 	for col_offset := -1; col_offset <= 1; col_offset++ {
	// 		possibleCollisionCoords = append(possibleCollisionCoords, []int{p.Y + row_offset, p.X + col_offset})
	// 	}
	// }
	newPos := p.Pos.Add(*NewVector2(int(math.Round(deltaX)), int(math.Round(deltaY))))
	for _, w := range walls {
		if newPos.Eq(*w.Pos) {
			return
		}
	}

	// Move
	p.Pos = &newPos
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
