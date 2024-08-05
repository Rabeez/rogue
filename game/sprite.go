package game

import (
	"image/color"
	"log"
	"slices"
	"time"

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

func (p *Sprite) Draw(panel *Panel) {
	// This method converts sprite grid coords to screen pixel coords for drawing
	op := &ebiten.DrawImageOptions{}
	pixelX := float64(p.Pos.X * TILE_SIZE)
	pixelY := float64(p.Pos.Y * TILE_SIZE)

	// Detect bounds and overflow drawing
	w, h := p.Img.Bounds().Dx(), p.Img.Bounds().Dy()
	if p.Pos.X+w > panel.Size.X || p.Pos.Y+h > panel.Size.Y {
		log.Fatalf("Overflow drawing detected\n")
	}

	// Account for panel offset
	op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
	op.GeoM.Translate(pixelX, pixelY)
	op.GeoM.Scale(2, 2)
	op.ColorScale.ScaleWithColor(p.color)
	panel.Screen.DrawImage(p.Img, op)
}

type Player struct {
	*Sprite
	speed       int
	health      int
	attackTimer *Timer
	gold        int
}

func NewPlayer(x, y int) *Player {
	// speed is in units of cells
	s := 1
	if s < 1 {
		log.Fatalf("Player speed has to be >= 1: %v", s)
	}
	return &Player{
		speed:       s,
		health:      5,
		attackTimer: NewTimer(time.Millisecond * 200),
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xff, 0xff, 0x00, 0xff},
			Img:   assets.PlayerSprite,
		},
	}
}

func (p *Player) Update(colliders *map[Vector2]bool, enemies *[]*Enemy, coins *[]*Coin) {
	p.attackTimer.Update()
	// All player keypresses
	deltaPos := NewVector2(0, 0)
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		deltaPos.Y = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		deltaPos.Y = p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		deltaPos.X = -p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		deltaPos.X = p.speed
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if p.attackTimer.IsReady() {
			p.attackTimer.Reset()
			// AOE attack in 3x3 around player with fixed 1 dmg
			var killedIdx []int
			for i, e := range *enemies {
				if d := p.Pos.ManDistance(*e.Pos); d <= 1 {
					e.health -= 1
					if e.health <= 0 {
						killedIdx = append(killedIdx, i)
					}
				}
			}
			for _, i := range killedIdx {
				*enemies = slices.Delete(*enemies, i, i+1)
			}
		}
	}
	newPos := p.Pos.Add(*deltaPos)

	// Collisions
	_, wallOverlap := (*colliders)[newPos]
	var enemyOverlap bool = false
	for _, e := range *enemies {
		enemyOverlap = newPos.Equals(*e.Pos)
		break
	}
	if !wallOverlap && !enemyOverlap {
		p.Pos = &newPos
	}

	// Pickups
	for i, c := range *coins {
		if c.Pos.Equals(*p.Pos) {
			p.gold += c.value
			*coins = slices.Delete(*coins, i, i+1)
			break
		}
	}

}

type Enemy struct {
	*Sprite
	aggroRadius   int
	health        int
	attackTimer   *Timer
	movementTimer *Timer
}

func NewEnemy(x, y int) *Enemy {
	ar := 3
	if ar < 1 {
		log.Fatalf("Enemy aggro radius has to be >= 1: %v", ar)
	}

	return &Enemy{
		aggroRadius:   ar,
		health:        2,
		attackTimer:   NewTimer(time.Second * 1),
		movementTimer: NewTimer(time.Millisecond * 200),
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xaa, 0x20, 0x20, 0xff},
			Img:   assets.EnemySprite,
		},
	}
}

func (e *Enemy) Update(p *Player, colliders *map[Vector2]bool) {
	e.movementTimer.Update()
	e.attackTimer.Update()

	if e.movementTimer.IsReady() {
		e.movementTimer.Reset()
		if d := e.Pos.ManDistance(*p.Pos); d <= e.aggroRadius {
			// fmt.Println(d)
			dir := p.Pos.Sub(*e.Pos).GridNormalize()
			newPos := e.Pos.Add(dir)
			// Collisions
			_, wallOverlap := (*colliders)[newPos]
			playerOverlap := newPos.Equals(*p.Pos)
			if !wallOverlap && !playerOverlap {
				e.Pos = &newPos
			}
		}
	}

	if e.attackTimer.IsReady() {
		e.attackTimer.Reset()
		// AOE attack in 3x3 around enemy with fixed 1 dmg
		if d := p.Pos.ManDistance(*e.Pos); d <= 1 {
			p.health -= 1
		}
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

type Coin struct {
	*Sprite
	value int
}

func NewCoin(x, y int, value int) *Coin {
	return &Coin{
		value: value,
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xaa, 0xaa, 0x00, 0xff},
			Img:   assets.CoinSprite,
		},
	}
}

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
