package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/Rabeez/rogue/assets"
)

type Sprite struct {
	X, Y float64
	Img  *ebiten.Image
}

type Player struct {
	*Sprite
}

func NewPlayer() *Player {
	return &Player{
		Sprite: &Sprite{
			X:   0,
			Y:   0,
			Img: assets.PlayerSprite,
		},
	}
}

type Enemy struct {
	*Sprite
}

func NewEnemy() *Enemy {
	return &Enemy{
		Sprite: &Sprite{
			X:   0,
			Y:   0,
			Img: assets.EnemySprite,
		},
	}
}

type Wall struct {
	*Sprite
}

func NewWall() *Wall {
	return &Wall{
		Sprite: &Sprite{
			X:   0,
			Y:   0,
			Img: assets.WallTSprite,
		},
	}
}
