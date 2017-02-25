package lgp

import "bitbucket.org/StephenPatrick/goevo/env"

type InitOptions struct {
	GenOpt           Options
	Env              *env.I
	Mem              *env.I
	Cross            LGPCrossover
	Act              []Action
	BaseActionWeight float64
	Fitness          FitnessFunc
	QuitEarly        int
}

func OptInit(opt interface{}) {
	iOpt := opt.(InitOptions)
	Init(
		iOpt.GenOpt,
		iOpt.Env,
		iOpt.Mem,
		iOpt.Cross,
		iOpt.Act,
		iOpt.BaseActionWeight,
		iOpt.Fitness,
		iOpt.QuitEarly,
	)
}

func Init(genOpt Options, e, m *env.I, cross LGPCrossover,
	act []Action, baseActionWeight float64, f FitnessFunc, qe int) {

	actions = act

	actionWeights = make([]float64, len(actions))
	for i := range actions {
		actionWeights[i] = baseActionWeight
	}
	ResetCumActionWeights()

	environment = e
	memory = m
	fitness = f
	gpOptions = genOpt
	crossover = cross
	quit_early = qe
}
