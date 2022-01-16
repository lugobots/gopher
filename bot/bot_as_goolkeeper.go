package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
	"math"
	"math/rand"
)

func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot, myState lugo4go.PlayerState) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	myGoal := field.GetTeamsGoal(b.side)

	if myState == lugo4go.HoldingTheBall {
		kickTarget := getTargetToRestartBall(b.side, snapshot)
		if kickTarget == nil {
			return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{}, "waiting players get closer"))
		}

		kickOrder, err := field.MakeOrderKickMaxSpeed(*snapshot.Ball, *kickTarget)
		if err != nil {
			return errors.Wrap(err, "error creating kicking order to pass from goalkeeper")
		}
		b.lastKickTurn = snapshot.Turn
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{kickOrder}, "restarting ball"))
	}

	stopOrder, err := field.MakeOrderMove(*me.Position, field.FieldCenter(), 0)
	if err != nil {
		return errors.Wrap(err, "error creating move order to stop goalkeeper")
	}

	if snapshot.Ball.Position.DistanceTo(*me.Position) > DistanceFar {
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{stopOrder, field.MakeOrderCatch()}, "just chilling"))
	}
	//1st - Based on the ball's speed in X axis, let find how long the ball would take to reach the coord X of the goal
	//2nd - find the nearest point the bal may reach at the goal
	//3rd - calculate the further point the keeper to be from that point to cover more area
	ballTarget, timeToReach, coming := findThreatenedSpot(*snapshot.GetBall(), myGoal)
	if !coming {
		optimumPoint := optimumWatchingPosition(myGoal, ballTarget, timeToReach)
		if optimumPoint.DistanceTo(*me.Position) < field.PlayerSize {
			return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{stopOrder, field.MakeOrderCatch()}, "I am ready"))
		}
		moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, optimumPoint)
		if err != nil {
			return errors.Wrap(err, "error creating move order to move goalkeeper")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "moving to a better position"))
	}

	distanceFromTarget := math.Abs(float64(ballTarget.Y - me.Position.Y))

	if distanceFromTarget < field.BallSize/2 {
		//the goalkeeper can already catch the ball!
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{stopOrder, field.MakeOrderCatch()}, "waiting to catch the ball"))
	}

	timeToCatch := int(distanceFromTarget / field.PlayerMaxSpeed)
	// if we do not have time to go running, let's JUMP!
	if timeToReach <= field.GoalKeeperJumpDuration && timeToCatch > timeToReach {
		idealSpeed := distanceFromTarget / field.GoalKeeperJumpDuration //we need ensure the jump won't be beyond the target
		jumpOrder, err := field.MakeOrderJump(*me.Position, ballTarget, idealSpeed)
		if err != nil {
			return errors.Wrap(err, "error creating jumping order for the goalkeeper")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{jumpOrder, field.MakeOrderCatch()}, "jumping to the success"))
	}

	//the keeper has time to catch the ball
	keeperSpeed := field.PlayerMaxSpeed
	if distanceFromTarget < field.PlayerMaxSpeed {
		keeperSpeed = distanceFromTarget //we do not want to pass through the ball target
	}
	moveOrder, err := field.MakeOrderMove(*me.Position, ballTarget, keeperSpeed)
	if err != nil {
		return errors.Wrap(err, "error creating move order to move the goalkeeper to the catching point")
	}

	return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "catching the ball!"))
}

// getTargetToRestartBall may return a nil point, which means there is no good target at this moment and the goalkeeper
// may wait.
func getTargetToRestartBall(mySide lugo.Team_Side, snapshot *lugo.GameSnapshot) *lugo.Point {
	me := field.GetPlayer(snapshot, mySide, field.GoalkeeperNumber)
	myTeam := field.GetTeam(snapshot, mySide).GetPlayers()
	opponentGoal := field.GetOpponentGoal(mySide)
	opponentSide := field.GetOpponentSide(mySide)
	opponentGoalKeeper := field.GetPlayer(snapshot, opponentSide, field.GoalkeeperNumber)
	opponentTeam := field.GetTeam(snapshot, opponentSide).GetPlayers()

	shouldIPass, bestCandidate := shouldIPass(me, snapshot.Ball.Position, opponentGoalKeeper.Position, opponentGoal, myTeam, opponentTeam)

	var target lugo.Point
	if shouldIPass < May {
		// if the goalkeeper cannot pass the ball yet, lest check if he can wait
		// if we still have at least N turns to hold it, let's hold it
		if snapshot.TurnsBallInGoalZone < field.BallTimeInGoalZone-3 {
			return nil
		}

		// we cannot wait, so let's kick randomly
		// @todo needs enhancement: we could find the place where there is no opponent player
		angleInDegrees := rand.Intn(90)
		if rand.Intn(2) == 0 {
			angleInDegrees *= -1
		}
		target := field.FieldCenter()
		randomDirection, _ := lugo.NewVector(*me.Position, target)
		randomDirection.AddAngleDegree(float64(angleInDegrees))
		target = randomDirection.TargetFrom(*me.Position)
	} else {
		target = *bestCandidate.Position
	}
	return &target
}

func findThreatenedSpot(ball lugo.Ball, myGoal field.Goal) (target lugo.Point, framesToReach int, coming bool) {
	ballSpeed := ball.Velocity.Speed
	ballXSpeed := ball.Velocity.Direction.Cos() * ballSpeed
	ballYSpeed := ball.Velocity.Direction.Sin() * ballSpeed

	if ball.Holder != nil {
		//let think what could happen if the ball was kicked now
		ballSpeed = field.BallMaxSpeed
		// if the ball wasn't kicked yet, the nearest point is the threatened
		target = NearestGoalPoint(ball, myGoal)

		ballKick, _ := lugo.NewVector(*ball.Position, target)
		ballXSpeed = ballKick.Cos() * ballSpeed
		ballYSpeed = ballKick.Sin() * ballSpeed
	}

	//S = So + V.T + (a/2).T^2
	//S: Goal X coord
	//So: Actual ball X coord
	//V: ballXSpeed
	//T: required
	//a: deceleration of the ball
	S := myGoal.Center.X
	So := ball.Position.X
	a := -field.BallDeceleration / 2
	// (a/2).T^2 +  V.T + (So - S)
	t1, t2, err := QuadraticResults(a, ballXSpeed, float64(So-S))
	if err != nil {
		return
	}
	realTimeToReach := t1 // truncating as integer because our time is calculated on frames
	if t1 <= 0 || (t2 > 0 && t2 < t1) {
		realTimeToReach = t2
	}

	if realTimeToReach < 0 {
		return
	}
	framesToReach = int(math.Ceil(realTimeToReach))

	// if the ball was kicked, let find the target based on its velocity
	//S: required
	//So: Actual ball Y coord
	//V: ballYSpeed
	//T:  "realTimeToReach"
	//a: deceleration of the ball
	coordY := float64(ball.Position.Y) + (ballYSpeed * realTimeToReach) + (a/2)*math.Pow(realTimeToReach, 2)

	target = lugo.Point{
		X: myGoal.Center.X,
		Y: int32(math.Round(coordY)),
	}

	coming = target.Y < field.GoalMaxY && target.Y > field.GoalMinY
	if ball.Holder != nil || ball.Velocity.Speed <= 0 {
		coming = false
	}
	return
}

func NearestGoalPoint(ball lugo.Ball, goalTarget field.Goal) lugo.Point {
	nearest := lugo.Point{
		X: goalTarget.Center.X,
		Y: ball.Position.Y,
	}
	if ball.Position.Y < field.GoalMinY {
		nearest.Y = goalTarget.BottomPole.Y + (field.BallSize / 2)
	} else if ball.Position.Y > field.GoalMaxY {
		nearest.Y = goalTarget.TopPole.Y - (field.BallSize / 2)
	}

	return nearest
}

// @todo it cam be enhanced: this function is not considering the player size, so the keeper could be further from the target sometimes
func optimumWatchingPosition(myGoal field.Goal, threatenedPoint lugo.Point, timeAvailable int) lugo.Point {
	jumpDistance := field.GoalKeeperJumpDuration * field.GoalKeeperJumpSpeed

	distanceFromCenter := myGoal.Center.DistanceTo(threatenedPoint)
	if jumpDistance > distanceFromCenter {
		distanceFromCenter -= jumpDistance
		timeAvailable -= field.GoalKeeperJumpDuration
	}

	if timeAvailable <= 1 { // too late!
		return threatenedPoint
	}

	//this is the time the keeper would take to reach the threatenedPoint if he started from the goal center
	timeNeededToReach := int(distanceFromCenter / field.PlayerMaxSpeed)
	if timeNeededToReach > timeAvailable {
		//the keeper is `lateTIme` frames late to reach the ball if the ball was kicked now
		lateTIme := timeNeededToReach - timeAvailable
		gapDistance := lateTIme * field.PlayerMaxSpeed
		//the gap is only in the Y axis, os it is easy to find the best point:

		optimumPoint := myGoal.Center
		if threatenedPoint.Y > myGoal.Center.Y { //above the center
			optimumPoint.Y += int32(gapDistance)
		} else {
			optimumPoint.Y -= int32(gapDistance)
		}
		return optimumPoint
	}

	//it's fine stay in the center, we have time enough to reach the threatenedPoint
	return myGoal.Center
}
