package location

import (
	"fmt"
	"math"
)

type Location struct {
	X float64
	Y float64
}

func (l *Location) String() string {
	return fmt.Sprintf("<%f,%f>", l.X, l.Y)
}

func (l *Location) Move(deltaX, deltaY float64) Location {
	return Location{l.X + deltaX, l.Y + deltaY}
}

func (l *Location) DistFrom(other Location) float64 {
	return math.Hypot(l.X-other.X, l.Y-other.Y)
}
