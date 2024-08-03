package game

import (
	"fmt"
	"log"
	"strings"

	"github.com/Rabeez/rogue/data"
	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	LevelNum  int
	Size      Vector2
	Player    *Player
	Enemies   []*Enemy
	Walls     []*Wall
	Colliders map[Vector2]bool
}

func makeLevelFromMatrix(mat [][]string) *Level {
	var p *Player
	var e []*Enemy
	var w []*Wall
	c := make(map[Vector2]bool)

	fmt.Printf("Generating level from %vX%v matrix\n", len(mat[0]), len(mat))
	sz := NewVector2(len(mat[0]), len(mat))

	for row, row_vals := range mat {
		for col, cell_val := range row_vals {
			vv := strings.Trim(cell_val, "\r\n ")
			if len(vv) == 0 {
				continue
			}
			switch vv {
			case "ep":
				p = NewPlayer(col, row)
			case "ee":
				e = append(e, NewEnemy(col, row))
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
			case "h":
				w = append(w, NewWall(col, row, Wall_Horz))
				c[*NewVector2(col, row)] = true
			case "v":
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

	// w := []*Wall{NewWall(100, 0, Wall_TopLeft), NewWall(100, 300, Wall_Top)}
	// p := NewPlayer(100, 100)
	// e := []*Enemy{NewEnemy(200, 100), NewEnemy(200, 300)}

	// return &Level{
	// 	LevelNum: level_num,
	// 	Player:   p,
	// 	Enemies:  e,
	// 	Walls:    w,
	// }

	l := makeLevelFromMatrix(data.TestLevel)
	return l
}

func (l *Level) Draw(screen *ebiten.Image) {
	screen.Fill(BACKGROUND_COLOR)

	for _, w := range l.Walls {
		w.Draw(screen)
	}
	for _, e := range l.Enemies {
		e.Draw(screen)
	}
	l.Player.Draw(screen)
}

func (l *Level) Update() error {
	l.Player.Update(&l.Colliders)

	return nil
}
