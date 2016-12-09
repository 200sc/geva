package goevo

import (
	"fmt"
	"goevo/population"
	"math/rand"
	"time"
)

func MakeDemes(demeCount int, members []population.Individual,
	s population.SelectionMethod, pair population.PairingMethod,
	in, out [][]float64, tests, elites, goal int, migration float64) population.DemeGroup {

	demeSize := len(members) / demeCount

	demes := make([]population.Population, demeCount)

	for i := 0; i < demeCount; i++ {
		demes[i] = population.Population{
			Members:      members[i*demeSize : (i+1)*demeSize],
			Size:         demeSize,
			Selection:    s,
			Pairing:      pair,
			FitnessTests: tests,
			TestInputs:   in,
			TestExpected: out,
			Elites:       elites,
			Fitnesses:    make([]int, demeSize),
			GoalFitness:  goal,
		}
	}

	return population.DemeGroup{
		Demes:           demes,
		MigrationChance: migration,
	}
}

func Seed() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RunDemeGroup(dg population.DemeGroup, numGens int) {
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
