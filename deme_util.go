package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/pop"
	"math/rand"
	"time"
)

func MakeDemes(demeCount int, members []pop.Individual,
	s []pop.SMethod, pair []pop.PMethod,
	in, out [][]float64, tests, goal int, elites alg.IntRange, migration float64) pop.DemeGroup {

	demeSize := len(members) / demeCount

	demes := make([]pop.Population, demeCount)

	for i := 0; i < demeCount; i++ {
		si := i % len(s)
		pi := i % len(pair)
		demes[i] = pop.Population{
			Members:      members[i*demeSize : (i+1)*demeSize],
			Size:         demeSize,
			Selection:    s[si],
			Pairing:      pair[pi],
			FitnessTests: tests,
			TestInputs:   in,
			TestExpected: out,
			Elites:       elites.Poll(),
			Fitnesses:    make([]int, demeSize),
			GoalFitness:  goal,
		}
	}

	return pop.DemeGroup{
		Demes:           demes,
		MigrationChance: migration,
	}
}

func Seed() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RunDemeGroup(dg pop.DemeGroup, numGens int) {
	for i := 0; i < numGens; i++ {
		fmt.Println("Gen", i+1)
		stopEarly := dg.NextGeneration()
		if i == numGens-1 || stopEarly {
			for _, p := range dg.Demes {
				w, _ := p.Weights(1.0)
				fmt.Println(w)
				maxWeight := 0.0
				maxIndex := 0
				for i, v := range w {
					if v > maxWeight {
						maxWeight = v
						maxIndex = i
					}
				}
				fmt.Println(p.Fitnesses)
				p.Members[maxIndex].Print()
			}
			fmt.Println("Generations taken: ", i+1)
			break
		}
	}
}
