package bot

import (
	"math/rand"
	"time"
)

type TeamState string

type Role string

func init() {
	rand.Seed(time.Now().UnixNano())
}

// IMPORTANT: all this constant sets below may be changed (see each set instructions). However, any change will
// affect the tactic defined in tactic.go file. So you must go there and adapt your tactics to your new settings.

// You however, may increase or decrease their values to change the precision of the Positioner.
// These values define how the field will be divided by the Positioner to create a field map.

// please update the tests if you include more states, or exclude some of them.
const (
	Initial       TeamState = "initial"
	UnderPressure TeamState = "under-pressure"
	Defensive     TeamState = "defensive"
	Neutral       TeamState = "neutral"
	Offensive     TeamState = "offensive"
	OnAttack      TeamState = "on-attack"
)
const (
	Defense Role = "defense"
	Middle  Role = "middle"
	Attack  Role = "attack"
)

type RegionCode struct {
	Col uint8
	Row uint8
}

type PlayerActionRegions map[TeamState]RegionCode

func DefinePlayerActionRegions(number uint32) PlayerActionRegions {
	return roleMap[number]
}

const (
	FieldGridCols = 10
	FieldGridRows = 8
)

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
