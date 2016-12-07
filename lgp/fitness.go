package lgp

type FitnessFunc func(gp *LGP, inputs, outputs [][]float64) int
