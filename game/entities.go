package game

import (
	"log"
	"math"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/teacat/noire"

	"github.com/Rabeez/rogue/assets"
)

type Sprite struct {
	// Sprite coords are in grid coord space
	Pos   *Vector2
	Img   *ebiten.Image
	color noire.Color
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
	op.ColorScale.ScaleWithColor(NoireToColor(s.color))
	panel.Screen.DrawImage(s.Img, op)
}

type Player struct {
	*Sprite
	speed                 int
	health                int
	attackTimer           *Timer
	gold                  int
	damageIndicatorTimers []*Timer
}

func NewPlayer(x, y int) *Player {
	// speed is in units of cells
	s := 1
	if s < 1 {
		log.Fatalf("Player speed has to be >= 1: %v", s)
	}
	return &Player{
		speed:                 s,
		health:                5,
		attackTimer:           NewTimer(time.Millisecond*200, true),
		gold:                  0,
		damageIndicatorTimers: []*Timer{},
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: noire.NewRGBA(0xff, 0xff, 0x00, 0xff),
			Img:   assets.PlayerSprite,
		},
	}
}

func (p *Player) Update(l *Level) {
	p.attackTimer.Update()
	// All player keypresses
	// TODO: implement RepeatingKeyPress to allow for holding keys
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
			// AOE attack in plus-shape around player with fixed 1 dmg
			for _, e := range l.Enemies {
				if d := p.Pos.ManDistance(*e.Pos); d <= 1 {
					// TODO: this should be set as 'payload' in timer? so damage number text uses correct value always
					dmg := 1
					e.health -= dmg
					// TODO: these values should be controlled fully within Enemy
					e.damageIndicatorTimer = NewTimer(time.Millisecond*300, false)
					if e.health <= 0 {
						e.isDying = true
					}
				}
			}

			// AOE attack in plus-shape around player
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

	// Update damage number timers
	for i, t := range p.damageIndicatorTimers {
		if t == nil {
			continue
		}
		t.Update()
		if t.IsReady() {
			p.damageIndicatorTimers[i] = nil
		}
	}
	// Count ended timers
	nils := 0
	for _, t := range p.damageIndicatorTimers {
		if t == nil {
			nils++
		}
	}
	// Only clear list if all timers have ended
	// This prevents poisitons being shifted between Draw calls
	if nils > 0 && nils == len(p.damageIndicatorTimers) {
		p.damageIndicatorTimers = []*Timer{}
	}

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
	if !p.attackTimer.IsReady() {
		op.ColorScale.ScaleWithColor(NoireToColor(p.color.Shade(p.attackTimer.CurrentProgress())))
	} else {
		op.ColorScale.ScaleWithColor(NoireToColor(p.color))
	}

	panel.Screen.DrawImage(p.Img, op)

	spriteW, spriteH := float64(p.Img.Bounds().Dx()), float64(p.Img.Bounds().Dy())
	xOffsets := *ArangeF(0, spriteW, len(p.damageIndicatorTimers))
	for i, t := range p.damageIndicatorTimers {
		if t == nil {
			continue
		}
		progress := t.CurrentProgress()
		txt := "1"
		op := &text.DrawOptions{}
		// Account for panel offset
		op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
		op.GeoM.Translate(pixelX, pixelY)
		// Distribute all numbers along sprite width
		op.GeoM.Translate(xOffsets[i], 0)
		// // random range
		// op.GeoM.Translate(RandRangeF(0, spriteW), 0)
		// TODO: center align text accounting for text width
		// textW, _ := text.Measure(txt, f, 0)
		// Animation
		op.GeoM.Translate(0, -(progress * spriteH / 2))
		op.GeoM.Scale(2, 2)

		PutText(panel, txt, op, NoireToColor(noire.NewRGBA(0xff, 0x00, 0x00, 0xff)), 12)
	}
}

type Enemy struct {
	*Sprite
	aggroRadius          int
	health               int
	attackTimer          *Timer
	movementTimer        *Timer
	isDying              bool
	isDead               bool
	deathTimer           *Timer
	damageIndicatorTimer *Timer
}

func NewEnemy(x, y int) *Enemy {
	ar := 3
	if ar < 1 {
		log.Fatalf("Enemy aggro radius has to be >= 1: %v", ar)
	}

	return &Enemy{
		aggroRadius:          ar,
		health:               2,
		attackTimer:          NewTimer(time.Second*1, true),
		movementTimer:        NewTimer(time.Millisecond*200, true),
		isDying:              false,
		isDead:               false,
		deathTimer:           nil,
		damageIndicatorTimer: nil,
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: noire.NewRGBA(0xaa, 0x20, 0x20, 0xff),
			Img:   assets.EnemySprite,
		},
	}
}

func (e *Enemy) Update(l *Level) {
	e.movementTimer.Update()
	e.attackTimer.Update()

	if e.isDying && e.deathTimer == nil {
		e.deathTimer = NewTimer(time.Millisecond*500*1, false)
	}

	if e.deathTimer != nil {
		e.deathTimer.Update()
		if e.deathTimer.IsReady() {
			e.isDead = true
		}
	}

	if e.damageIndicatorTimer != nil {
		e.damageIndicatorTimer.Update()
		if e.damageIndicatorTimer.IsReady() {
			e.damageIndicatorTimer = nil
		}
	}

	if !e.isDead && e.movementTimer.IsReady() {
		e.movementTimer.Reset()
		if d := e.Pos.ManDistance(*l.Player.Pos); d <= e.aggroRadius {
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

	if !e.isDead && e.attackTimer.IsReady() {
		// AOE attack in plus-shape around enemy with fixed 1 dmg
		if d := l.Player.Pos.ManDistance(*e.Pos); d <= 1 {
			e.attackTimer.Reset()
			// TODO: this should be set as 'payload' in timer? so damage number text uses correct value always
			dmg := 1
			l.Player.health -= dmg
			// TODO: these values should be controlled fully within Player
			l.Player.damageIndicatorTimers = append(l.Player.damageIndicatorTimers, NewTimer(time.Millisecond*300, false))
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

	// Dynamic color based on attack readiness or dying state
	op.ColorScale.ScaleWithColor(NoireToColor(e.color))
	if e.isDying {
		if !e.deathTimer.IsReady() {
			p := math.Max(0.5, e.deathTimer.CurrentProgress())
			op.ColorScale.ScaleAlpha(1 - float32(p))
		}
	} else {
		if !e.attackTimer.IsReady() {
			p := e.attackTimer.CurrentProgress()
			t := math.Pow(math.Min(p, 1-p), 2)
			op.ColorScale.ScaleWithColor(NoireToColor(e.color.Shade(t)))
		} else {
			op.ColorScale.ScaleWithColor(NoireToColor(e.color))
		}
	}

	panel.Screen.DrawImage(e.Img, op)

	if e.damageIndicatorTimer != nil {
		progress := e.damageIndicatorTimer.CurrentProgress()
		txt := "1"
		op := &text.DrawOptions{}
		// Account for panel offset
		op.GeoM.Translate(float64(panel.Corner.X), float64(panel.Corner.Y))
		op.GeoM.Translate(pixelX, pixelY)
		spriteW, spriteH := float64(e.Img.Bounds().Dx()), float64(e.Img.Bounds().Dy())
		// Center align at middle of sprite
		op.GeoM.Translate(spriteW/2.0, 0)
		// TODO: center align text accounting for text width
		// textW, _ := text.Measure(txt, f, 0)
		// Animation
		op.GeoM.Translate(0, -(progress * spriteH / 2))
		op.GeoM.Scale(2, 2)

		PutText(panel, txt, op, NoireToColor(noire.NewRGBA(0xff, 0x00, 0x00, 0xff)), 12)
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
		color: noire.NewRGBA(0xff, 0xff, 0xff, 0xff),
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
