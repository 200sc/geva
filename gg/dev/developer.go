package dev

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/200sc/geva/env"
	"github.com/200sc/geva/pop"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
)

type Dev interface {
	Mechanic() *Mechanic
	SetFitness(int)
	pop.Individual
}

type Base struct {
	InitSize          intrange.Range
	GoalSize          intrange.Range
	GoalDistance      intrange.Range
	EnvSize           intrange.Range
	EnvVal            floatrange.Range
	ActionCount       intrange.Range
	ActionTypes       []ActionType
	ActionTypeWeights []float64
	ActionStrengths   []floatrange.Range
	//PassiveRatio      float64

	// Fitness is set by controlling instances
	fitness int

	// Crossover and Mutation methods set by generators
	DevMutation
	Cross Crossover
}

func (d *Base) String() string {
	s := "Base Dev\n"
	s += "InitSize:" + iRangeString(d.InitSize)
	s += " GoalSize:" + iRangeString(d.GoalSize)
	s += " GoalDistance:" + iRangeString(d.GoalDistance)
	s += " EnvSize:" + iRangeString(d.EnvSize)
	s += " EnvVal:" + fRangeString(d.EnvVal)
	s += " ActionCount:" + iRangeString(d.ActionCount)
	return s
}

func iRangeString(ir intrange.Range) string {
	i1 := ir.Percentile(0)
	i2 := ir.Percentile(1)
	if i2 < i1 {
		i1, i2 = i2, i1
	}
	return "Range: " + strconv.Itoa(i1) + ", " + strconv.Itoa(i1)
}

func fRangeString(fr floatrange.Range) string {
	f1 := strconv.FormatFloat(fr.Percentile(0), 'E', -1, 64)
	f2 := strconv.FormatFloat(fr.Percentile(1), 'E', -1, 64)
	return "Range: " + f1 + ", " + f2 + "\n"
}

func (d *Base) SetFitness(f int) {
	d.fitness = f
}

func (d *Base) Fitness(input, expected [][]float64) int {
	return d.fitness
}

func (d *Base) Mutate() {
	if d == nil {
		fmt.Println("d is nil")
	}
	if d.DevMutation == nil {
		fmt.Println("d's mutation is nil")
	}
	d.DevMutation.Mutate(d)
}

func (d *Base) Crossover(other pop.Individual) pop.Individual {
	d2 := other.(*Base)
	return d.Cross.Crossover(d, d2)
}

func (d *Base) CanCrossover(other pop.Individual) bool {
	_, ok := other.(*Base)
	return ok
}

func (d *Base) Print() {
	fmt.Println("A developer")
}

func Default() *Base {
	dev := new(Base)
	dev.InitSize = intrange.NewConstant(1)
	dev.GoalSize = intrange.NewConstant(1)
	dev.GoalDistance = intrange.NewConstant(5)
	dev.EnvSize = intrange.NewConstant(5)
	dev.EnvVal = floatrange.NewConstant(0)
	dev.ActionCount = intrange.NewConstant(5)
	dev.ActionTypes = BaseActionTypes
	dev.ActionStrengths = BaseActionStrength
	dev.ActionTypeWeights = BaseActionWeights
	return dev
}

// Developers produce Mechanics

func (d *Base) Mechanic() *Mechanic {
	gg := new(Mechanic)

	e := env.NewF(d.EnvSize.Poll(), 0.0)

	// For each environment variable,
	// Generate a number of actions based on ActionCount
	// that modify that variable, the type of which chosen through roulette
	// search on the cumulative weights of ActionTypeWeights,
	// resolved from ActionTypes to Actions using some strength
	// based on ActionStrengths.
	actions := make([]func(), 0)
	cum := alg.CumulativeWeights(d.ActionTypeWeights)
	for i := 0; i < len(*e); i++ {
		actCt := d.ActionCount.Poll()
		for j := 0; j < actCt; j++ {
			k := alg.CumWeightedChooseOne(cum)
			a := actionMut(d.ActionTypes[k](d.ActionStrengths[k].Poll()), e.GetP(i))
			actions = append(actions, a)
		}
	}

	// scramble actions
	// for i := 0; i < len(actions); i++ {
	// 	j := rand.Intn(len(actions))
	// 	actions[i], actions[j] = actions[j], actions[i]
	// }
	gg.Actions = actions

	if len(gg.Actions) == 0 {
		fmt.Println("No actions")
		fmt.Println("Environment length:", len(*e), d.EnvSize.Percentile(0))
		fmt.Println("ActionCount min:", d.ActionCount.Percentile(0))
	}

	// splitIndex := int(math.Ceil(float64(len(actions)) * d.PassiveRatio))

	// gg.Actions = actions[0:splitIndex]
	// gg.Passives = actions[splitIndex : len(actions)+1]

	gg.Init = make(map[int]float64)
	// Choose some number of variables to initialize
	// at game start
	l := d.InitSize.Poll()
	for i := 0; i < l; i++ {
		j := rand.Intn(len(*e))
		gg.Init[j] = d.EnvVal.Poll()
	}

	gg.Goal = make(map[int]float64)

	// Simulate some actions on the environment . . .
	e.SetMap(gg.Init)

	l = d.GoalDistance.Poll()
	for i := 0; i < l; i++ {
		//Perform some action
		gg.Actions[rand.Intn(len(gg.Actions))]()
		//Perform all passive actions
		// for j, p := range gg.Passives {
		// 	p()
		// }
	}

	// . . . and pull random elements from said environment
	// to determine the goal state.
	l = d.GoalSize.Poll()
	for i := 0; i < l; i++ {
		j := rand.Intn(len(*e))
		gg.Goal[j] = e.Get(j)
	}

	e.SetAll(0.0)
	gg.Environment = e

	gg.MechFitness = (*Mechanic).FitnessElems

	return gg
}
