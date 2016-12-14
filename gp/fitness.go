package gp

import (
	"math"
	"time"
)

type FitnessFunc func(gp *GP, inputs, outputs [][]float64) int

// An example fitness function which treats
// the output as a environment to compare
// a modified environment by the GP to.
func EnvFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		Eval(g.First)
		fitness += g.Env.Diff(outputs[i])
	}
	return fitness
}

func MatchEnvFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		Eval(g.First)
		fitness += g.Env.MatchDiff(outputs[i])
	}
	return fitness
}

func MatchMemFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		Eval(g.First)
		fitness += g.Mem.MatchDiff(outputs[i])
	}
	return fitness
}

// An example fitness which treats the
// output of the GP as a value to compare
// against the single expected output
func OutputFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		out := Eval(g.First)
		fitness += int(math.Abs(float64(out - int(outputs[i][0]))))
	}
	if fitness < 0 {
		fitness = math.MaxInt32
	}
	return fitness
}

func Mem0Fitness(g *GP, inputs, outputs [][]float64) int {
	if g.Mem == nil {
		panic("Mem0Fitness used on GPs without memory")
	}
	fitness := 1
	for i, envDiff := range inputs {
		g.Env = environment.New(envDiff)
		Eval(g.First)
		fitness += int(math.Abs(float64(*(*g.Mem)[0]) - outputs[i][0]))
	}
	if fitness < 0 {
		fitness = math.MaxInt32
	}
	return fitness
}

func ComplexityFitness(f FitnessFunc, mod float64) FitnessFunc {
	return func(g *GP, inputs, outputs [][]float64) int {
		i := f(g, inputs, outputs)
		i += int(math.Floor(mod * float64(g.Nodes)))
		if i < 0 {
			i = math.MaxInt32
		}
		return i
	}
}

func TimeFitness(f FitnessFunc, threshold int, timeLimit int) FitnessFunc {
	return func(g *GP, inputs, outputs [][]float64) int {
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
