package bot

import (
	"context"

	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/pkg/errors"
)

type Bot struct {
	mapper        field.Mapper
	Role          Role
	number        uint32
	side          proto.Team_Side
	actionRegions PlayerActionRegions
	log           lugo4go.Logger

	lastKickTurn uint32
}

func NewBot(logger lugo4go.Logger, side proto.Team_Side, number uint32) *Bot {
	fieldMapper, _ := field.NewMapper(FieldGridCols, FieldGridRows, side)

	me := Bot{
		mapper: fieldMapper,
		Role:   DefineRole(number),
		number: number,
		side:   side,
		log:    logger,
	}
	if number != field.GoalkeeperNumber {
		me.actionRegions = DefinePlayerActionRegions(number)
	}
	return &me
}

func (b *Bot) OnDisputing(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	ballPosition := snapshot.GetBall().GetPosition()

	ballRegion, err := b.mapper.GetPointRegion(ballPosition)
	if err != nil {
		return errors.Wrap(err, "failed to find the ball region")
	}

	myRegion, err := b.mapper.GetPointRegion(me.GetPosition())
	if err != nil {
		return errors.Wrap(err, "failed to find my region")
	}

	// if the ball is max 2 blocks away from me, I will move toward the ball
	if !isNear(myRegion, ballRegion) {
		return b.holdPosition(ctx, sender, snapshot)
	}

	moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *ballPosition)
	if err != nil {
		return errors.Wrap(err, "error creating move order")
	}

	// we can ALWAYS try to catch the ball if we are not holding it
	catchOrder := field.MakeOrderCatch()

	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder, catchOrder}, "trying to catch the ball"))

}

func (b *Bot) OnDefending(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	ballPosition := snapshot.GetBall().GetPosition()

	ballRegion, err := b.mapper.GetPointRegion(ballPosition)
	if err != nil {
		return errors.Wrap(err, "failed to find the ball region")
	}

	myRegion, err := b.mapper.GetPointRegion(me.GetPosition())
	if err != nil {
		return errors.Wrap(err, "failed to find my region")
	}

	// if the ball is max 2 blocks away from me, I will move toward the ball
	if !isNear(myRegion, ballRegion) {
		return b.holdPosition(ctx, sender, snapshot)
	}

	moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *ballPosition)
	if err != nil {
		return errors.Wrap(err, "error creating move order")
	}

	// we can ALWAYS try to catch the ball if we are not holding it
	catchOrder := field.MakeOrderCatch()

	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder, catchOrder}, "trying to catch the ball"))
}

func (b *Bot) OnHolding(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	goal := field.GetOpponentGoal(b.side)

	goalRegion, err := b.mapper.GetPointRegion(&goal.Center)
	if err != nil {
		return errors.Wrap(err, "failed to find the ball region")
	}

	myRegion, err := b.mapper.GetPointRegion(me.GetPosition())
	if err != nil {
		return errors.Wrap(err, "failed to find my region")
	}

	var order proto.PlayerOrder
	// if we are near to the goal, let's kick it!
	if !isNear(myRegion, goalRegion) {
		return b.holdPosition(ctx, sender, snapshot)

	}

	order, err = field.MakeOrderKickMaxSpeed(*snapshot.GetBall(), goal.Center)
	if err != nil {
		return errors.Wrap(err, "failed to create kick order")
	}

	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{order}, "attack!"))
}

func (b *Bot) OnSupporting(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	teammatePosition := snapshot.GetBall().GetHolder().GetPosition()

	teammateRegion, err := b.mapper.GetPointRegion(teammatePosition)
	if err != nil {
		return errors.Wrap(err, "failed to find the teammate region")
	}

	myRegion, err := b.mapper.GetPointRegion(me.GetPosition())
	if err != nil {
		return errors.Wrap(err, "failed to find my region")
	}

	// if the player mate is max 2 blocks away from me, I will help
	if !isNear(myRegion, teammateRegion) {
		return b.holdPosition(ctx, sender, snapshot)
	}

	moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *teammatePosition)
	if err != nil {
		return errors.Wrap(err, "error creating move order")
	}

	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder}, "trying to catch the ball"))
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot, myState lugo4go.PlayerState) error {
	me := field.GetPlayer(snapshot, b.side, b.number)

	if myState == lugo4go.HoldingTheBall {
		order, err := field.MakeOrderKickMaxSpeed(*snapshot.GetBall(), field.FieldCenter())
		if err != nil {
			return errors.Wrap(err, "failed to create kick order")
		}

		return processServerResp(sender.Send(ctx, []proto.PlayerOrder{order}, "kick the ball"))
	}

	ballPosition := snapshot.GetBall().GetPosition()

	moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *ballPosition)
	if err != nil {
		return errors.Wrap(err, "error creating move order")
	}

	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder, field.MakeOrderCatch()}, "trying to catch the ball"))
}

// Aux methods

func (b *Bot) MyInitialPosition() *proto.Point {
	iniRegion := b.actionRegions[Initial]
	// we may ignore this error because if it is not a valid region we will notice during the development
	region, _ := b.mapper.GetRegion(iniRegion.Col, iniRegion.Row)
	return region.Center()
}

func (b *Bot) myActionRegion(teamState TeamState) field.Region {
	r, _ := b.mapper.GetRegion(b.actionRegions[teamState].Col, b.actionRegions[teamState].Row)
	return r
}

func (b *Bot) holdPosition(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	teamState := Neutral

	ballRegion, _ := b.mapper.GetPointRegion(snapshot.Ball.Position)
	if snapshot.GetShotClock() != nil {
		teamState, _ = DetermineTeamState(ballRegion, b.side, snapshot.GetShotClock().GetTeamSide())
	}

	actionRegion := b.myActionRegion(teamState)

	speed := field.PlayerMaxSpeed
	msg := "moving to my region"
	if me.Position.DistanceTo(*actionRegion.Center()) < field.PlayerSize {
		speed = 0
		msg = "Holding position"
	}

	moveOrder, err := field.MakeOrderMove(*me.Position, *actionRegion.Center(), speed)
	if err != nil {

		return errors.Wrap(err, "error creating move order to return to action region")
	}
	return processServerResp(sender.Send(ctx, []proto.PlayerOrder{moveOrder, field.MakeOrderCatch()}, msg))
}
