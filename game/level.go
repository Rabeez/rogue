package game

import (
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/Rabeez/rogue/data"
)

type Level struct {
	LevelNum int
	Size     Vector2
	Player   *Player
	Enemies  []*Enemy
	Coins    []*Coin
	Chests   []*Chest
	// TODO: only keep map for wall data and loop over it for drawing. Separate array is unnecesasry
	Walls []*Wall
	// TODO: map Vec2 -> array of sprite interface(?)
	Colliders map[Vector2]bool
}

func makeLevelFromMatrix(mat [][]string) *Level {
	var p *Player
	var e []*Enemy
	var coins []*Coin
	var chests []*Chest
	var w []*Wall
	c := make(map[Vector2]bool)

	sz := NewVector2(len(mat[0]), len(mat))
	log.Printf("Generating level from %vx%v matrix\n", sz.X, sz.Y)

	for row, row_vals := range mat {
		for col, cell_val := range row_vals {
			vv := strings.Trim(cell_val, "\r\n ")
			if len(vv) == 0 {
				continue
			}
			vv_check := vv[:2]
			switch vv_check {
			case "ep":
				p = NewPlayer(col, row)
			case "ee":
				e = append(e, NewEnemy(col, row))
			case "ch":
				chests = append(chests, NewChest(col, row))
			case "ec":
				coin_val, err := strconv.Atoi(vv[2:])
				if err != nil {
					panic(err)
				}
				coins = append(coins, NewCoin(col, row, coin_val))
			case "tl":
				w = append(w, NewWall(col, row, Wall_TopLeft))
				c[*NewVector2(col, row)] = true
			case "tt":
				w = append(w, NewWall(col, row, Wall_TopT))
				c[*NewVector2(col, row)] = true
			case "tr":
				w = append(w, NewWall(col, row, Wall_TopRight))
				c[*NewVector2(col, row)] = true
			case "lt":
				w = append(w, NewWall(col, row, Wall_LeftT))
				c[*NewVector2(col, row)] = true
			case "m":
				w = append(w, NewWall(col, row, Wall_Middle))
				c[*NewVector2(col, row)] = true
			case "rt":
				w = append(w, NewWall(col, row, Wall_RightT))
				c[*NewVector2(col, row)] = true
			case "bl":
				w = append(w, NewWall(col, row, Wall_LowerLeft))
				c[*NewVector2(col, row)] = true
			case "bt":
				w = append(w, NewWall(col, row, Wall_LowerT))
				c[*NewVector2(col, row)] = true
			case "br":
				w = append(w, NewWall(col, row, Wall_LowerRight))
				c[*NewVector2(col, row)] = true
			case "hw":
				w = append(w, NewWall(col, row, Wall_Horz))
				c[*NewVector2(col, row)] = true
			case "vw":
				w = append(w, NewWall(col, row, Wall_Vert))
				c[*NewVector2(col, row)] = true
			default:
				log.Fatalf("Invalid wall label found in level: '%v'", cell_val)
			}
		}
	}

	l := Level{
		LevelNum:  -1,
		Size:      *sz,
		Player:    p,
		Enemies:   e,
		Coins:     coins,
		Chests:    chests,
		Walls:     w,
		Colliders: c,
	}
	return &l
}

func NewLevel(level_num int) *Level {
	// TODO: level gen
	// separate method for generation? maybe different variants with different algos
	// algo will run in grid coords and return appropriate arrays to fill Level struct
	// will need a camera type in game w/ 'zoom' level
	// camera center will be on player and zoom will be num of tiles radius around
	// will shift visible slice on player movement
	// only render visible slice
	// but simulaiton will run on whole grid BBOX
	// optim:
	// keep bigger buffer around radius as off-screen render for smoother movement
	// move the visible slice in separate goroutine in background for smoother movement without hitches?
	// extra:
	// (zoom will change on player speed -> potions etc or on enemy hit zoom in shake?)

	// TODO: seed rng

	l := makeLevelFromMatrix(data.TestLevel)
	return l
}

func (l *Level) Draw(panel *Panel) {
	panel.Screen.Fill(BACKGROUND_COLOR)

	for _, w := range l.Walls {
		w.Draw(panel)
	}
	for _, e := range l.Enemies {
		e.Draw(panel)
	}
	for _, c := range l.Chests {
		c.Draw(panel)
	}
	for _, c := range l.Coins {
		c.Draw(panel)
	}
	l.Player.Draw(panel)
}

func (l *Level) Update() error {
	l.Player.Update(l)

	var deadThisTick []int
	for i, e := range l.Enemies {
		e.Update(l)
		if e.isDead {
			deadThisTick = append(deadThisTick, i)
		}
	}
	for _, i := range deadThisTick {
		l.Enemies = slices.Delete(l.Enemies, i, i+1)
	}

	return nil
}
