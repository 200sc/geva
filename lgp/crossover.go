package lgp

type LGPCrossover interface {
	Crossover(a, b *LGP) *LGP
}

type PointCrossover struct {
	NumPoints int
}

func (pc PointCrossover) Crossover(a, b *LGP) *LGP {

}

type UniformCrossover struct {
	ChosenProportion float64
}

func (uc UniformCrossover) Crossover(a, b *LGP) *LGP {

}
