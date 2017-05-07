package eda

import (
	"math"
	"math/rand"
)

// I'm not sure how we want to model failure in the onemax problem.
// given a distribution, should we say it is as wrong as all of its
// elements are far away from 1, or should we evaluate each element
// as a probability, and say each element is wrong if a random number
// does not fall under the float?
//
// These two ideas are represented here separately.

func OnemaxABS(m Model) int {
	e := m.ToEnv()
	diff := 0.0
	for i, f := range *env {
		diff += math.Abs(*(*env)[i] - 1)
	}
	return int(diff)
}

func OnemaxChance(m Model) int {
	e := m.ToEnv()
	diff := 0
	for i, f := range *env {
		if rand.Float64() > *f {
			diff++
		}
	}
	return diff
}
