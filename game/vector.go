package game

type Vector2 struct {
	X, Y int
}

func NewVector2(x, y int) *Vector2 {
	return &Vector2{
		X: x,
		Y: y,
	}
}

func (v Vector2) Eq(other Vector2) bool {
	return v.X == other.X && v.Y == other.Y
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}
