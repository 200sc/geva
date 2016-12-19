package neural

type InitOptions struct {
	Ngo     NetworkGenerationOptions
	Cross   NeuralCrossover
	Fitness FitnessFunc
}

func OptInit(opt interface{}) {
	iOpt := opt.(InitOptions)
	Init(
		iOpt.Ngo,
		iOpt.Cross,
		iOpt.Fitness,
	)
}

func Init(newNgo NetworkGenerationOptions, newCrossover NeuralCrossover, f FitnessFunc) {
	ngo = newNgo
	crossover = newCrossover
	fitness = f
}
