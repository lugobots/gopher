package bot

import (
	"github.com/lugobots/lugo4go/v2/pkg/field"
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

type FuzzyScale int

const (
	MustNot FuzzyScale = iota
	ShouldNot
	May
	Should
	Must
)

const (
	DistanceBeside = field.FieldWidth / 10
	DistanceNear   = field.FieldWidth / 8
	DistanceFar    = field.FieldWidth / 6
	DistanceTooFar = field.FieldWidth / 4
)

type RegionCode struct {
	Col uint8
	Row uint8
}

type PlayerActionRegions map[TeamState]RegionCode

func DefinePlayerActionRegions(number uint32) PlayerActionRegions {
	return roleMap[number]
}
