package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
)

func (b *Bot) OnDefending(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	if ShouldIDisputeForTheBall(b.mapper, b.number, me.Position, snapshot.Ball.Position, field.GetTeam(snapshot, b.side).Players) {
		speed, target := FindBestPointInterceptBall(snapshot.GetBall(), me)
		moveOrder, err := field.MakeOrderMove(*me.Position, *target, speed)
		if err != nil {
			return errors.Wrap(err, "error creating move order")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "trying to take the ball"))
	}
	ballRegion, _ := b.mapper.GetPointRegion(snapshot.Ball.Position)
	teamState, _ := DetermineTeamState(ballRegion, b.side, snapshot.GetShotClock().GetTeamSide())
	if b.Role == Defense && teamState == UnderPressure {
		myGOal := field.GetTeamsGoal(b.side)
		opponentTrajectoryToGoal, _ := lugo.NewVector(*snapshot.Ball.Holder.Position, myGOal.Center)
		myTrajectoryToGoal, _ := lugo.NewVector(*me.Position, myGOal.Center)

		opponentTrajectoryToGoal.AngleWith(myTrajectoryToGoal)
		if opponentTrajectoryToGoal.AngleDegrees() > 0 {
			myTrajectoryToGoal.Perpendicular()
		} else {
			myTrajectoryToGoal.Invert().Perpendicular()
		}
		targetDirection := myTrajectoryToGoal.TargetFrom(*me.Position)
		moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, targetDirection)
		if err != nil {
			return errors.Wrap(err, "error creating move order")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "trying to take the ball"))

	}

	return b.holdPosition(ctx, sender, snapshot)
}
