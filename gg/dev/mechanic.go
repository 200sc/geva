package dev

import (
	"math"

	"github.com/200sc/geva/env"
)

type Mechanic struct {
	Actions []func()
	//Passives []func()
	//Continuous bool
	//PassiveTypes     []ActionType
	//PassiveStrengths []float64
	//
	Environment *env.F
	Init        map[int]float64
	Goal        map[int]float64
	MechFitness
}

func (mch *Mechanic) Reset() {

}

type MechFitness func(*Mechanic) int

func (mch *Mechanic) FitnessElems() int {
	// opt 1.
	// How many goal elements are in position
	fitness := 1
	for i, v := range mch.Goal {
		if mch.Environment.Get(i) != v {
			fitness++
		}
	}
	return fitness
}

func (mch *Mechanic) FitnessAbs() int {
	// opt 2.
	// Absolute distance of all goal elements from position
	fitness := 1.0
	for i, v := range mch.Goal {
		v2 := mch.Environment.Get(i)
		fitness += math.Abs(v - v2)
	}
	return int(fitness) + 1
}
