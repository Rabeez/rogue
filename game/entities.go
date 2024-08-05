package game

import (
	"image/color"
	"log"
	"math"
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

func (s *Sprite) Draw(panel *Panel) {
	// This method converts sprite grid coords to screen pixel coords for drawing
	op := &ebiten.DrawImageOptions{}
	pixelX := float64(s.Pos.X * TILE_SIZE)
	pixelY := float64(s.Pos.Y * TILE_SIZE)

	// Detect bounds and overflow drawing
	w, h := s.Img.Bounds().Dx(), s.Img.Bounds().Dy()
	if s.Pos.X+w > panel.Size.X || s.Pos.Y+h > panel.Size.Y {
		log.Fatalf("Overflow drawing detected\n")
	}

	// Account for panel offset
	op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
	op.GeoM.Translate(pixelX, pixelY)
	op.GeoM.Scale(2, 2)
	op.ColorScale.ScaleWithColor(s.color)
	panel.Screen.DrawImage(s.Img, op)
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
		attackTimer: NewTimer(time.Millisecond*200, true),
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xff, 0xff, 0x00, 0xff},
			Img:   assets.PlayerSprite,
		},
	}
}

func (p *Player) Update(l *Level) {
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
			for i, e := range l.Enemies {
				if d := p.Pos.ManDistance(*e.Pos); d <= 1 {
					e.health -= 1
					// TODO: start a text anim here for damage numbers
					if e.health <= 0 {
						killedIdx = append(killedIdx, i)
					}
				}
			}
			// TODO: have enemy update it's sprite while in "dying state"
			// article has example with plane changing shape
			for _, i := range killedIdx {
				l.Enemies = slices.Delete(l.Enemies, i, i+1)
			}

			// AOE attack in 3x3 around player
			var brokenIdx []int
			for i, c := range l.Chests {
				if d := p.Pos.ManDistance(*c.Pos); d <= 1 {
					brokenIdx = append(brokenIdx, i)
				}
			}
			for _, i := range brokenIdx {
				// TODO: return loottable here and directly modify level?
				l.Chests[i].Open(&l.Coins)
				l.Chests = slices.Delete(l.Chests, i, i+1)
			}
		}
	}
	newPos := p.Pos.Add(*deltaPos)

	// Collisions
	_, wallOverlap := l.Colliders[newPos]
	var enemyOverlap bool = false
	for _, e := range l.Enemies {
		enemyOverlap = newPos.Equals(*e.Pos)
		break
	}
	if !wallOverlap && !enemyOverlap {
		p.Pos = &newPos
	}

	// Pickups
	for i, c := range l.Coins {
		if c.Pos.Equals(*p.Pos) {
			p.gold += c.value
			l.Coins = slices.Delete(l.Coins, i, i+1)
			break
		}
	}
}

func (p *Player) Draw(panel *Panel) {
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

	// Dynamic color based on attack readiness
	// TODO: color animation while taking damage
	op.ColorScale.ScaleWithColor(p.color)
	if !p.attackTimer.IsReady() {
		op.ColorScale.ScaleAlpha(float32(math.Max(0.5, p.attackTimer.CurrentProgress())))
	}

	panel.Screen.DrawImage(p.Img, op)
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
		attackTimer:   NewTimer(time.Second*1, true),
		movementTimer: NewTimer(time.Millisecond*200, true),
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xaa, 0x20, 0x20, 0xff},
			Img:   assets.EnemySprite,
		},
	}
}

func (e *Enemy) Update(l *Level) {
	e.movementTimer.Update()
	e.attackTimer.Update()

	if e.movementTimer.IsReady() {
		e.movementTimer.Reset()
		if d := e.Pos.ManDistance(*l.Player.Pos); d <= e.aggroRadius {
			// fmt.Println(d)
			dir := l.Player.Pos.Sub(*e.Pos).GridNormalize()
			newPos := e.Pos.Add(dir)
			// Collisions
			_, wallOverlap := l.Colliders[newPos]
			playerOverlap := newPos.Equals(*l.Player.Pos)
			if !wallOverlap && !playerOverlap {
				e.Pos = &newPos
			}
		}
	}

	if e.attackTimer.IsReady() {
		// AOE attack in 3x3 around enemy with fixed 1 dmg
		if d := l.Player.Pos.ManDistance(*e.Pos); d <= 1 {
			e.attackTimer.Reset()
			// TODO: start a text anim here for damage numbers
			l.Player.health -= 1
		}
	}
}

func (e *Enemy) Draw(panel *Panel) {
	// This method converts sprite grid coords to screen pixel coords for drawing
	op := &ebiten.DrawImageOptions{}
	pixelX := float64(e.Pos.X * TILE_SIZE)
	pixelY := float64(e.Pos.Y * TILE_SIZE)

	// Detect bounds and overflow drawing
	w, h := e.Img.Bounds().Dx(), e.Img.Bounds().Dy()
	if e.Pos.X+w > panel.Size.X || e.Pos.Y+h > panel.Size.Y {
		log.Fatalf("Overflow drawing detected\n")
	}

	// Account for panel offset
	op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
	op.GeoM.Translate(pixelX, pixelY)
	op.GeoM.Scale(2, 2)

	// Dynamic color based on attack readiness
	// TODO: color animation while taking damage
	op.ColorScale.ScaleWithColor(e.color)
	if !e.attackTimer.IsReady() {
		op.ColorScale.ScaleAlpha(float32(math.Max(0.5, e.attackTimer.CurrentProgress())))
	}

	panel.Screen.DrawImage(e.Img, op)
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
