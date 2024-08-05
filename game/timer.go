package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Timer struct {
	currentTicks int
	targetTicks  int
}

func NewTimer(d time.Duration, startReady bool) *Timer {
	full := int(d.Milliseconds()) * ebiten.TPS() / 1000
	starter := 0
	if startReady {
		starter = full
	}
	return &Timer{
		currentTicks: starter,
		targetTicks:  full,
	}
}

func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

func (t *Timer) IsReady() bool {
	return t.currentTicks >= t.targetTicks
}

func (t *Timer) Reset() {
	t.currentTicks = 0
}

func (t *Timer) CurrentProgress() float64 {
	return math.Min(float64(t.currentTicks)/float64(t.targetTicks), 1.0)
}
