package bot

import (
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"math"
)

func DistanceBetweenRegions(a, b field.Region) float64 {
	return math.Hypot(
		math.Abs(float64(a.Col()-b.Col())),
		math.Abs(float64(b.Col()-b.Col())),
	)
}
