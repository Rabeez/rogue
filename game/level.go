package game

import "github.com/hajimehoshi/ebiten/v2"

type Level struct {
	LevelNum int
	Player   *Player
	Enemies  []*Enemy
	Walls    []*Wall
}

func NewLevel(level_num int) *Level {

	// TODO: level gen
	// separate method for generation? maybe different variants with different algos
	// output should be a 2D array (maybe TileType enum?)
	// will need a camera type in game w/ 'zoom' level
	// camera center will be on player and zoom will be num of tiles radius around
	// will shift visible slice on player movement
	// only render visible slice
	// but simulaiton will run on whole 2d array
	// optim:
	// keep bigger buffer around radius as off-screen render for smoother movement
	// move the visible slice in separate goroutine in background for smoother movement without hitches?
	// extra:
	// (zoom will change on player speed -> potions etc or on enemy hit zoom in shake?)

	w := []*Wall{NewWall(100, 0, Wall_TopLeft), NewWall(100, 300, Wall_Top)}
	p := NewPlayer(100, 100)
	e := []*Enemy{NewEnemy(200, 100), NewEnemy(200, 300)}

	return &Level{
		LevelNum: level_num,
		Player:   p,
		Enemies:  e,
		Walls:    w,
	}
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
	l.Player.Update()

	return nil
}
