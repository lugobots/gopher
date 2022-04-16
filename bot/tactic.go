package bot

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

const (
	FieldGridCols = 10
	FieldGridRows = 8
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

var roleMap = map[uint32]PlayerActionRegions{
	// starting from 2 because the number goalkeeper has PlayerActionRegions
	2: {
		Initial:       {1, 2},
		UnderPressure: {0, 3},
		Defensive:     {1, 3},
		Neutral:       {2, 3},
		Offensive:     {2, 3},
		OnAttack:      {5, 3},
	},
	3: {
		Initial:       {1, 5},
		UnderPressure: {0, 4},
		Defensive:     {1, 4},
		Neutral:       {2, 4},
		Offensive:     {2, 4},
		OnAttack:      {5, 4},
	},
	4: {
		Initial:       {2, 6},
		UnderPressure: {1, 6},
		Defensive:     {2, 6},
		Neutral:       {3, 6},
		Offensive:     {5, 5},
		OnAttack:      {6, 5},
	},
	5: {
		Initial:       {2, 1},
		UnderPressure: {1, 1},
		Defensive:     {2, 1},
		Neutral:       {3, 1},
		Offensive:     {5, 2},
		OnAttack:      {6, 2},
	},
	6: {
		Initial:       {3, 6},
		UnderPressure: {1, 5},
		Defensive:     {2, 5},
		Neutral:       {4, 5},
		Offensive:     {8, 7},
		OnAttack:      {9, 6},
	},
	7: {
		Initial:       {3, 1},
		UnderPressure: {1, 2},
		Defensive:     {2, 2},
		Neutral:       {4, 2},
		Offensive:     {8, 0},
		OnAttack:      {9, 1},
	},
	8: {
		Initial:       {3, 3},
		UnderPressure: {1, 3},
		Defensive:     {3, 3},
		Neutral:       {4, 3},
		Offensive:     {6, 3},
		OnAttack:      {7, 3},
	},
	9: {
		Initial:       {3, 4},
		UnderPressure: {1, 4},
		Defensive:     {3, 4},
		Neutral:       {4, 4},
		Offensive:     {6, 4},
		OnAttack:      {7, 4},
	},
	10: {
		Initial:       {4, 5},
		UnderPressure: {2, 5},
		Defensive:     {4, 5},
		Neutral:       {5, 6},
		Offensive:     {7, 5},
		OnAttack:      {8, 5},
	},
	11: {
		Initial:       {4, 2},
		UnderPressure: {2, 2},
		Defensive:     {4, 2},
		Neutral:       {5, 1},
		Offensive:     {7, 2},
		OnAttack:      {9, 2},
	},
}
