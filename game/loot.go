package game

import (
	"github.com/Rabeez/rogue/assets"
	"github.com/teacat/noire"
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
			color: noire.NewRGBA(0xaa, 0xaa, 0x00, 0xff),
			Img:   assets.CoinSprite,
		},
	}
}

type LootTable struct {
	gold int
}

func NewLootTable() *LootTable {
	gold := RandRange(2, 10)
	return &LootTable{
		gold: gold,
	}
}

var possibleDropOffsets = []Vector2{
	*NewVector2(0, 0),
	*NewVector2(0, -1),
	*NewVector2(0, 1),
	*NewVector2(1, 0),
	*NewVector2(-1, 0),
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
			color: noire.NewRGBA(0xaa, 0xaa, 0x00, 0xff),
			Img:   assets.ChestSprite,
		},
	}
}

func (l *Chest) Open(coins *[]*Coin) {
	dropOffsetIdx := 0
	{
		dropLocation := NewVector2(l.Pos.X+possibleDropOffsets[dropOffsetIdx].X, l.Pos.Y+possibleDropOffsets[dropOffsetIdx].Y)
		*coins = append(*coins, NewCoin(int(dropLocation.X), int(dropLocation.Y), l.loot.gold))
		dropOffsetIdx++
	}
}
