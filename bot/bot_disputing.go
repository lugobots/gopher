package bot

import (
	"context"
	"github.com/lugobots/arena/units"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
	"math"
)

func processServerResp(resp *lugo.OrderResponse, err error) error {
	if err != nil {
		return errors.Wrapf(err, "error sending orders")
	}
	if resp.Code == lugo.OrderResponse_SUCCESS {
		return errors.Errorf("server responded a non-success code: %s", resp.GetCode().Descriptor())
	}
	return nil
}

func (b *Bot) OnDisputing(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	if ShouldIDisputeForTheBall(b.mapper, b.number, me.Position, snapshot.Ball.Position, field.GetTeam(snapshot, b.side).Players) {
		speed, target := FindBestPointInterceptBall(snapshot.GetBall(), me)
		moveOrder, err := field.MakeOrderMove(*me.Position, *target, speed)
		if err != nil {
			return errors.Wrap(err, "error creating move order")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "disputing for the ball"))
	}
	ballRegion, _ := b.mapper.GetPointRegion(snapshot.Ball.Position)
	teamState, _ := DetermineTeamState(ballRegion, b.side, b.side)
	actionRegion := b.myActionRegion(teamState)
	if currentRegion, _ := b.mapper.GetPointRegion(me.Position); currentRegion != b.myActionRegion(teamState) {
		moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *actionRegion.Center())
		if err != nil {
			return errors.Wrap(err, "error creating move order to return to action region")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "moving to my region"))
	}

	moveOrder, _ := field.MakeOrderMoveMaxSpeed(*me.Position, field.FieldCenter())
	return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder}, "Holding position"))
}

func ShouldIDisputeForTheBall(mapper field.Mapper, botNumber uint32, botPosition, ballPosition *lugo.Point, teamMates []*lugo.Player) bool {
	ballRegion, _ := mapper.GetPointRegion(ballPosition)
	botRegion, _ := mapper.GetPointRegion(botPosition)
	if DistanceBetweenRegions(botRegion, ballRegion) < 2 {
		return true
	}
	myDistance := ballPosition.DistanceTo(*ballPosition)
	playerCloser := 0
	for _, teamMate := range teamMates {
		if teamMate.Number != botNumber && teamMate.Position.DistanceTo(*ballPosition) < myDistance {
			playerCloser++
			if playerCloser > 1 { // are there more than on player closer to the ball than me?
				return false
			}
		}
	}
	return true
}

func FindBestPointInterceptBall(ball *lugo.Ball, player *lugo.Player) (speed float64, target *lugo.Point) {
	if ball.Velocity.Speed == 0 {
		return field.PlayerMaxSpeed, ball.Position
	} else {
		// @todo needs enhancement: there are math formulas to find the sweet spot
		calcBallPos := func(frame int) *lugo.Point {
			//S = So + VT + (aT^2)/2
			V := ball.Velocity.Speed
			T := float64(frame)
			a := -units.BallDeceleration
			distance := V*T + (a*math.Pow(T, 2))/2
			if distance <= 0 {
				return nil
			}
			vectorToBal, _ := ball.Velocity.Direction.Copy().SetLength(distance)
			ballTarget := vectorToBal.TargetFrom(*ball.Position)
			return &ballTarget
		}
		frames := 1
		lastBallPosition := ball.Position
		for {
			ballLocation := calcBallPos(frames)
			if ballLocation == nil {
				break
			}
			minDistanceToTouch := ballLocation.DistanceTo(*player.Position) - ((units.BallSize + units.PlayerSize) / 2)

			if minDistanceToTouch <= float64(units.PlayerMaxSpeed*frames) {
				if frames > 1 {
					return units.PlayerMaxSpeed, ballLocation
				} else {
					return player.Position.DistanceTo(*ballLocation), ballLocation
				}
			}
			lastBallPosition = ballLocation
			frames++
		}
		return units.PlayerMaxSpeed, lastBallPosition
	}
}
