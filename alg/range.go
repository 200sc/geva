package alg

import (
	"math/rand"
)

// IntRange represents the ability
// to poll a struct and return an integer,
// distributed over some range dependant
// on the implementing struct.
type IntRange interface {
	Poll() int
}

// LinearIntRange polls on a linear scale
// between a minimum and a maximum
type LinearIntRange struct {
	Min, Max int
}

// Poll returns an integer distributed
// between lir.Min and lir.Max
func (lir LinearIntRange) Poll() int {
	return rand.Intn(lir.Max-lir.Min) + lir.Min
}

// Constant implements IntRange as a poll
// which always returns the same integer.
type Constant int

// Poll returns c cast to an int
func (c Constant) Poll() int {
	return int(c)
}
