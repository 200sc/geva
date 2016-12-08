package lgp

import (
	"math"
	"time"
)

type FitnessFunc func(gp *LGP, inputs, outputs [][]float64) int

// An example fitness function which treats
// the output as a environment to compare
// a modified environment by the GP to.
func EnvFitness(g *LGP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		g.Run()
		fitness += g.Env.Diff(outputs[i])
	}
	return fitness
}

func Mem0Fitness(g *LGP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		g.Run()
		fitness += int(math.Abs(float64(*(*g.Mem)[0]) - outputs[i][0]))
	}
	if fitness < 0 {
		fitness = math.MaxInt32
	}
	return fitness
}

func ComplexityFitness(f FitnessFunc, mod float64) FitnessFunc {
	return func(g *LGP, inputs, outputs [][]float64) int {
		i := f(g, inputs, outputs)
		i += int(math.Floor(mod * float64(len(g.Instructions))))
		if i < 0 {
			i = math.MaxInt32
		}
		return i
	}
}

func TimeFitness(f FitnessFunc, threshold int, timeLimit int) FitnessFunc {
	return func(g *LGP, inputs, outputs [][]float64) int {
		t1 := time.Now()
		i := f(g, inputs, outputs)
		t2 := time.Now().Sub(t1)

		if i <= threshold {
			t3 := int(t2 / time.Second)
			if t3 < timeLimit {
				i -= int(math.Floor(float64(threshold) * (float64(t3) / float64(timeLimit))))
			}
		}
		i += threshold
		if i < 0 {
			i = math.MaxInt32
		}
		return i
	}
}
