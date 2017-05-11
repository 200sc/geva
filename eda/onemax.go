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

// MaxABS is an abstracted form of Onemax for targetting any
// floating point value
func MaxABS(t float64) func(m Model) int {
	return func(m Model) int {
		e := m.ToEnv()
		diff := 0.0
		for _, f := range *e {
			diff += math.Abs(*f - t)
		}
		return int(diff)
	}
}

// OnemaxABS is a fitness function which returns the absolute difference
// in an environment from an environment of all ones.
var OnemaxABS = MaxABS(1)

// OnemaxChance is a fitness function which rolls rng on every value in
// an environment and returns the number of values which were not rolled
// under.
func OnemaxChance(m Model) int {
	e := m.ToEnv()
	diff := 0
	for _, f := range *e {
		if rand.Float64() > *f {
			diff++
		}
	}
	return diff
}

// Preliminary experimental results seem to show that OnemaxABS learns faster
// than OnemaxChance, however that could just be due to it being a faster
// fitness function to run (no rng rolls needed)
