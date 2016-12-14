package alg

import (
	"math/rand"
)

type IntRange interface {
	Poll() int
}

type LinearIntRange struct {
	Min, Max int
}

func (lir LinearIntRange) Poll() int {
	return rand.Intn(lir.Max-lir.Min) + lir.Min
}
