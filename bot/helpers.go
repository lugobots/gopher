package bot

import (
	"fmt"
	"math"

	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/pkg/errors"
)

func DefineRole(number uint32) Role {
	// starting from 2 because the number goalkeeper has no role
	switch number {
	case 2, 3, 4, 5:
		return Defense
	case 6, 7, 8, 9:
		return Middle
	case 10, 11:
		return Attack
	}
	return ""
}

func DetermineTeamState(ballRegion field.Region, myTeamSide, possession proto.Team_Side) (s TeamState, e error) {
	regionCol := ballRegion.Col()
	if possession == myTeamSide {
		switch regionCol {
		case 5, 6, 7, 8, 9:
			return OnAttack, nil
		case 2, 3, 4:
			return Offensive, nil
		case 0, 1:
			return Neutral, nil
		}

	} else {
		switch regionCol {
		case 9:
			return Defensive, nil
		case 6, 7, 8:
			return Defensive, nil
		case 3, 4, 5:
			return Defensive, nil
		case 0, 1, 2:
			return UnderPressure, nil
		}
		//return Offensive, nil
	}
	return "", fmt.Errorf("unknown team state for ball in %d col, tor possion with %s", regionCol, possession)
}

func isNear(a, b field.Region) bool {
	const minDist = 2
	colDist := a.Col() - b.Col()
	rowDist := a.Row() - b.Row()
	return math.Hypot(float64(colDist), float64(rowDist)) <= minDist
}

func processServerResp(resp *proto.OrderResponse, err error) error {
	if err != nil {
		return errors.Wrapf(err, "error sending orders")
	}
	if resp.Code != proto.OrderResponse_SUCCESS {
		return errors.Errorf("server responded a non-success code: %s", resp.Code.String())
	}
	return nil
}
