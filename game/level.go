package game

import "github.com/hajimehoshi/ebiten/v2"

type Level struct {
	LevelNum int
	Player   *Player
	Enemies  []*Enemy
	Walls    []*Wall
}

func NewLevel(level_num int) *Level {

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
