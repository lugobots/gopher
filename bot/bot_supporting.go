package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
)

const numberOfAssistsPlayers = 3
const assistPlayerDistance = DistanceBeside

func (b *Bot) OnSupporting(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	opponentSide := field.GetOpponentSide(b.side)
	myTeam := field.GetTeam(snapshot, b.side).GetPlayers()
	opponentTeam := field.GetTeam(snapshot, opponentSide).GetPlayers()

	shouldIAssist, _ := ShouldIAssist(me.Position, snapshot.Ball.Position, b.number, snapshot.Ball.Holder.Number, myTeam)
	if shouldIAssist {
		bestSpot := FindSpotsToAssist(snapshot.Ball.Holder.Position, me.Position, b.side, myTeam, opponentTeam)
		speed := field.PlayerMaxSpeed
		msg := "on my way"
		if me.Position.DistanceTo(*bestSpot) < field.PlayerSize {
			speed = 0
			msg = "I am here"
		}
		moveOrder, err := field.MakeOrderMove(*me.Position, *bestSpot, speed)
		if err != nil {
			return errors.Wrap(err, "error creating moving order to assist")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder}, msg))
	}
	return b.holdPosition(ctx, sender, snapshot)
}

func ShouldIAssist(myPosition, ballPosition *lugo.Point, myNumber, holderNumber uint32, myTeam []*lugo.Player) (bool, []*lugo.Player) {
	playerCloser := 0
	myDistance := myPosition.DistanceTo(*ballPosition)

	assistingPlayers := make([]*lugo.Player, 0, numberOfAssistsPlayers)
	// should have at least 3 supporters in the perimeters around the ball holder
	for _, teammate := range myTeam {
		if teammate.Number != holderNumber && // the holder cannot help himself
			teammate.Number != myNumber && // I won't count to myself
			teammate.Position.DistanceTo(*ballPosition) < myDistance {
			assistingPlayers = append(assistingPlayers, teammate)
			playerCloser++
			if playerCloser >= numberOfAssistsPlayers { // are there more than two player closer to the ball than me?
				return false, assistingPlayers
			}
		}
	}
	return true, assistingPlayers
}

func FindSpotsToAssist(ballHolderPosition, botPosition *lugo.Point, teamSide lugo.Team_Side, assistPlayer, opponentTeam []*lugo.Player) *lugo.Point {
	backPositionY := ballHolderPosition.Y - assistPlayerDistance
	if teamSide == lugo.Team_AWAY {
		backPositionY = ballHolderPosition.Y + assistPlayerDistance
	}

	perfectPositionBack := &lugo.Point{
		X: ballHolderPosition.X,
		Y: backPositionY,
	}
	perfectPositionLateralA := &lugo.Point{
		X: ballHolderPosition.X,
		Y: ballHolderPosition.Y - assistPlayerDistance,
	}
	perfectPositionLateralB := &lugo.Point{
		X: ballHolderPosition.X,
		Y: ballHolderPosition.Y + assistPlayerDistance,
	}

	positions := []*lugo.Point{perfectPositionBack, perfectPositionLateralA, perfectPositionLateralB}
	for i, originalPos := range positions {
		obstaclesBack, _ := findOpponentsOnMyRoute(originalPos, ballHolderPosition, field.PlayerSize, opponentTeam)
		if len(obstaclesBack) > 0 {
			v, err := lugo.NewVector(*botPosition, *obstaclesBack[0].position)
			if err != nil {
				if obstaclesBack[0].distanceFromTrajectory > 0 {
					v.Y = -v.Y
				} else {
					v.X = -v.X
				}
				_, _ = v.SetLength(field.PlayerSize)
				fixedPosition := v.TargetFrom(*originalPos)
				positions[i] = &fixedPosition
			}
		}
	}

	// @todo needs enhancement: the closest spot is not always the most efficient because we must consider the other players
	closestDistance := float64(field.FieldWidth)
	var closestPosition lugo.Point
	for _, position := range positions {
		d := position.DistanceTo(*botPosition)
		if d < closestDistance {
			closestDistance = d
			closestPosition = *position
		}
	}

	return &closestPosition
}
