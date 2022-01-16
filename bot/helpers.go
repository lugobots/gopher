package bot

import (
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
	"math"
	"sort"
)

func DistanceBetweenRegions(a, b field.Region) float64 {
	return math.Hypot(
		math.Abs(float64(a.Col()-b.Col())),
		math.Abs(float64(b.Col()-b.Col())),
	)
}

type obstacleDetails struct {
	position               *lugo.Point
	distanceFromTrajectory float64
}

func findOpponentsOnMyRoute(origin, target *lugo.Point, margin float64, opponentTeam []*lugo.Player) ([]obstacleDetails, error) {
	var obstacles []obstacleDetails

	minX := int32(math.Min(float64(origin.X), float64(target.X)))
	maxX := int32(math.Max(float64(origin.X), float64(target.X)))
	minY := int32(math.Min(float64(origin.Y), float64(target.Y)))
	maxY := int32(math.Max(float64(origin.Y), float64(target.Y)))

	for _, opponent := range opponentTeam {
		pointX := opponent.Position.X
		pointY := opponent.Position.Y

		x1 := origin.X
		y1 := origin.Y
		x2 := target.X
		y2 := target.Y

		// formula found at "Line defined by two points"  https://en.wikipedia.org/wiki/Distance_from_a_point_to_a_line
		distanceToTrajectory := (float64((x2-x1)*(y1-pointY)) - float64((x1-pointX)*(y2-y1))) /
			math.Sqrt(math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2))

		if math.Abs(distanceToTrajectory) <= margin &&
			((opponent.Position.X > minX && opponent.Position.X < maxX) ||
				(opponent.Position.Y > minY && opponent.Position.Y < maxY)) {
			obstacles = append(obstacles, obstacleDetails{
				position:               opponent.Position,
				distanceFromTrajectory: distanceToTrajectory,
			})
		}
	}
	sort.Slice(obstacles, func(i, j int) bool {
		return math.Abs(obstacles[i].distanceFromTrajectory) < math.Abs(obstacles[j].distanceFromTrajectory)
	})

	return obstacles, nil
}

// QuadraticResults resolves a quadratic function returning the x1 and x2
func QuadraticResults(a, b, c float64) (float64, float64, error) {
	if a == 0 {
		return 0, 0, errors.New("a cannot be zero")
	}
	// delta: B^2 -4.A.C
	delta := math.Pow(b, 2) - 4*a*c
	// quadratic formula: -b +/- sqrt(delta)/2a
	t1 := (-b + math.Sqrt(delta)) / (2 * a)
	t2 := (-b - math.Sqrt(delta)) / (2 * a)
	return t1, t2, nil
}
