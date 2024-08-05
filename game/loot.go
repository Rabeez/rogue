package game

import (
	"image/color"

	"github.com/Rabeez/rogue/assets"
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

type LootTable struct {
	gold int
}

func NewLootTable() *LootTable {
	gold := randRange(2, 10)
	return &LootTable{
		gold: gold,
	}
}

type Chest struct {
	*Sprite
	loot *LootTable
}

func NewChest(x, y int) *Chest {
	return &Chest{
		loot: NewLootTable(),
		Sprite: &Sprite{
			Pos:   NewVector2(x, y),
			color: color.RGBA{0xaa, 0xaa, 0x00, 0xff},
			Img:   assets.ChestSprite,
		},
	}
}

func (l *Chest) Open(coins *[]*Coin) {
	*coins = append(*coins, NewCoin(l.Pos.X, l.Pos.Y, l.loot.gold))
}
