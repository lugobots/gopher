package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/pkg/errors"
)

const numberOfAssistsPlayers = 3

func (b *Bot) OnSupporting(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	opponentSide := field.GetOpponentSide(b.side)
	myTeam := field.GetTeam(snapshot, b.side).GetPlayers()
	opponentTeam := field.GetTeam(snapshot, opponentSide).GetPlayers()

	shouldIAssist, _ := ShouldIAssist(me.Position, snapshot.Ball.Position, b.number, snapshot.Ball.Holder.Number, myTeam)
	if shouldIAssist {
		bestSpot := FindSpotsToAssist(snapshot.Ball.Holder.Position, me.Position, b.mapper, opponentTeam)
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
		return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder}, msg))
	}
	return b.holdPosition(ctx, sender, snapshot)
}

func ShouldIAssist(myPosition, ballPosition *proto.Point, myNumber, holderNumber uint32, myTeam []*proto.Player) (bool, []*proto.Player) {
	playerCloser := 0
	myDistance := myPosition.DistanceTo(*ballPosition)

	assistingPlayers := make([]*proto.Player, 0, numberOfAssistsPlayers)
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

func FindSpotsToAssist(ballHolderPosition, botPosition *proto.Point, mapper field.Mapper, opponentTeam []*proto.Player) *proto.Point {
	holderRegion, _ := mapper.GetPointRegion(ballHolderPosition)

	perfectPositionBack := holderRegion.Back().Center()
	perfectPositionLateralA := holderRegion.Left().Front().Center()
	perfectPositionLateralB := holderRegion.Right().Front().Center()

	positions := []*proto.Point{perfectPositionBack, perfectPositionLateralA, perfectPositionLateralB}
	for i, originalPos := range positions {
		obstaclesBack, _ := findOpponentsOnMyRoute(originalPos, ballHolderPosition, field.PlayerSize, opponentTeam)
		if len(obstaclesBack) > 0 {
			v, err := proto.NewVector(*botPosition, *obstaclesBack[0].position)
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

	// @todo needs enhancement should check if another player is closer to that point
	// @todo needs enhancement: the closest spot is not always the most efficient because we must consider the other players
	closestDistance := float64(field.FieldWidth)
	var closestPosition proto.Point
	for _, position := range positions {
		d := position.DistanceTo(*botPosition)
		if d < closestDistance {
			closestDistance = d
			closestPosition = *position
		}
	}

	return &closestPosition
}
