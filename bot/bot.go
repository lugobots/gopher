package bot

import (
	"context"
	"fmt"
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

func (b *Bot) OnDefending(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	ballRegion, _ := b.mapper.GetPointRegion(snapshot.Ball.Position)
	teamState, _ := DetermineTeamState(ballRegion, b.side, snapshot.GetShotClock().GetTeamSide())
	actionRegion := b.myActionRegion(teamState)
	if currentRegion, _ := b.mapper.GetPointRegion(me.Position); !currentRegion.Eq(b.myActionRegion(teamState)) {
		moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *actionRegion.Center())
		if err != nil {
			return errors.Wrap(err, "error creating move order to return to action region")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder},
			fmt.Sprintf("deffending moving to my region %dx%d from %dx%d",
				actionRegion.Col(), actionRegion.Row(),
				currentRegion.Col(), currentRegion.Row(),
			)))
	}
	return nil
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot, state lugo4go.PlayerState) error {
	// nothing
	return nil
}

func (b *Bot) holdPosition(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *lugo.GameSnapshot) error {
	me := field.GetPlayer(snapshot, b.side, b.number)
	ballRegion, _ := b.mapper.GetPointRegion(snapshot.Ball.Position)
	teamState, _ := DetermineTeamState(ballRegion, b.side, b.side)
	actionRegion := b.myActionRegion(teamState)
	if currentRegion, _ := b.mapper.GetPointRegion(me.Position); currentRegion != b.myActionRegion(teamState) {
		moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, *actionRegion.Center())
		if err != nil {
			return errors.Wrap(err, "error creating move order to return to action region")
		}
		return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder}, "moving to my region"))
	}

	moveOrder, _ := field.MakeOrderMoveMaxSpeed(*me.Position, field.FieldCenter())
	return processServerResp(sender.Send(ctx, []lugo.PlayerOrder{moveOrder}, "Holding position"))
}
