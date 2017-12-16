package dev

import (
	"fmt"
	"math"
	"strconv"

	"github.com/200sc/geva/env"
	"github.com/200sc/geva/unique"
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
	mch.Environment.SetAll(0)
	for i, v := range mch.Init {
		mch.Environment.Set(i, v)
	}
}

type MechNode interface {
	MechanicDistance(*Mechanic) float64
}

func (mch *Mechanic) Distance(n unique.Node) (float64, bool) {
	if mch2, ok := n.(MechNode); ok {
		return mch2.MechanicDistance(mch), true
	}
	return 0, false
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

func (mch *Mechanic) String() string {
	s := "Mechanic:\n"
	s += "ActionCount: " + strconv.Itoa(len(mch.Actions))
	s += " Environment: " + mch.Environment.String()
	s += " Init: " + fmt.Sprint(mch.Init)
	s += " Goal: " + fmt.Sprint(mch.Goal)
	return s
}
