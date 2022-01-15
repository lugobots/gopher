package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/pkg/errors"
)

type Bot struct {
	mapper        field.Mapper
	Role          Role
	number        uint32
	side          lugo.Team_Side
	actionRegions PlayerActionRegions
	log           lugo4go.Logger

	lastKickTurn uint32
}

const afterKickingWaitingTime = uint32(8)

func NewBot(logger lugo4go.Logger, side lugo.Team_Side, number uint32) *Bot {
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

func (b *Bot) MyInitialPosition() *lugo.Point {
	iniRegion := b.actionRegions[Initial]
	// we may ignore this error because if it is not a valid region we will notice during the development
	region, _ := b.mapper.GetRegion(iniRegion.Col, iniRegion.Row)
	return region.Center()
}

func (b *Bot) myActionRegion(teamState TeamState) field.Region {
	r, _ := b.mapper.GetRegion(b.actionRegions[teamState].Col, b.actionRegions[teamState].Row)
	return r
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot, state lugo4go.PlayerState) error {
	// nothing
	return nil
}

func (b *Bot) holdPosition(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	teamState := Neutral

	//interval := (snapshot.Turn / 100) % 5
	//if interval < 1 {
	//	teamState = UnderPressure
	//} else if interval < 2 {
	//	teamState = Defensive
	//} else if interval < 3 {
	//	teamState = Neutral
	//} else if interval < 4 {
	//	teamState = Offensive
	//} else {
	//	teamState = OnAttack
	//}

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
	return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder, field.MakeOrderCatch()}, msg))
}
