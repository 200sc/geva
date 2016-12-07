package lgp

type LGPCrossover interface {
	Crossover(a, b *LGP) *LGP
}

type PointCrossover struct {
	NumPoints int
}

type UniformCrossover struct {
	ChosenProportion float64
}
