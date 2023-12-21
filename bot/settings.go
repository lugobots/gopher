package bot

import (
	"github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/mapper"
)

func GetPlayerTacticRegion(inspector lugo4go.SnapshotInspector, fieldMap mapper.Mapper, playerNumber int) mapper.Region {
	ballRegion, _ := fieldMap.GetPointRegion(inspector.GetBall().GetPosition())
	regionCol := ballRegion.Col()

	teamState := "defensive"
	if regionCol > 11 {
		teamState = "offensive"
	} else if regionCol > 6 {
		teamState = "neutral"
	}

	pos := DefaultTacticPositions[teamState][playerNumber]
	reg, _ := fieldMap.GetRegion(pos.Col, pos.Row)
	return reg
}

var DefaultTacticPositions = map[string]map[int]struct {
	Col int
	Row int
}{

	"initial": {
		2:  {Col: 1, Row: 2},
		3:  {Col: 1, Row: 3},
		4:  {Col: 1, Row: 4},
		5:  {Col: 1, Row: 5},
		6:  {Col: 4, Row: 1},
		7:  {Col: 4, Row: 3},
		8:  {Col: 4, Row: 4},
		9:  {Col: 4, Row: 6},
		10: {Col: 6, Row: 3},
		11: {Col: 6, Row: 4},
	},
	"defensive": {
		2:  {Col: 1, Row: 2},
		3:  {Col: 1, Row: 3},
		4:  {Col: 1, Row: 4},
		5:  {Col: 1, Row: 5},
		6:  {Col: 4, Row: 1},
		7:  {Col: 4, Row: 3},
		8:  {Col: 4, Row: 4},
		9:  {Col: 4, Row: 6},
		10: {Col: 6, Row: 3},
		11: {Col: 6, Row: 4},
	},
	"neutral": {
		2:  {Col: 3, Row: 1},
		3:  {Col: 3, Row: 3},
		4:  {Col: 3, Row: 4},
		5:  {Col: 3, Row: 6},
		6:  {Col: 6, Row: 1},
		7:  {Col: 6, Row: 3},
		8:  {Col: 6, Row: 4},
		9:  {Col: 6, Row: 6},
		10: {Col: 10, Row: 2},
		11: {Col: 10, Row: 5},
	},
	"offensive": {
		2:  {Col: 5, Row: 1},
		3:  {Col: 4, Row: 3},
		4:  {Col: 4, Row: 4},
		5:  {Col: 5, Row: 6},
		6:  {Col: 11, Row: 1},
		7:  {Col: 9, Row: 2},
		8:  {Col: 9, Row: 5},
		9:  {Col: 11, Row: 6},
		10: {Col: 13, Row: 2},
		11: {Col: 13, Row: 5},
	},
}
