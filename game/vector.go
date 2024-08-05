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

func (v Vector2) Equals(other Vector2) bool {
	return v.X == other.X && v.Y == other.Y
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector2) ManDistance(other Vector2) int {
	return abs(v.X-other.X) + abs(v.Y-other.Y)
}

func (v Vector2) GridNormalize() Vector2 {
	if abs(v.X) >= abs(v.Y) {
		return *NewVector2(sign(v.X), 0)
	} else {
		return *NewVector2(0, sign(v.Y))
	}
}
