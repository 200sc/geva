package gp

import "github.com/200sc/geva/env"

type InitOptions struct {
	GenOpt           Options
	Env              *env.I
	Cross            GPCrossover
	Act              [][]Action
	BaseActionWeight float64
	Fitness          FitnessFunc
	StorageCount     int
	StorageWeight    float64
}

func OptInit(opt interface{}) {
	iOpt := opt.(InitOptions)
	Init(
		iOpt.GenOpt,
		iOpt.Env,
		iOpt.Cross,
		iOpt.Act,
		iOpt.BaseActionWeight,
		iOpt.Fitness,
	)
	if iOpt.StorageCount != 0 {
		AddStorage(iOpt.StorageCount, iOpt.StorageWeight)
	}
}

func Init(genOpt Options, e *env.I, cross GPCrossover,
	act [][]Action, baseActionWeight float64, f FitnessFunc) {

	environment = e
	actions = act
	actionWeights = make([][]float64, len(actions))
	for i, tier := range actions {
		actionWeights[i] = make([]float64, len(tier))
		for j := range tier {
			actionWeights[i][j] = baseActionWeight
		}
	}
	cumZeroActionWeights = CalculateCumulativeActionWeights(0)
	cumActionWeights = CalculateCumulativeActionWeights(1, 2, 3)
	fitness = f
	gpOptions = genOpt
	crossover = cross
}
