package population

import (
	"math"
	"math/rand"
)

type DemeGroup struct {
	Demes []Population
	// Migration Chance is the average number
	// of individuals that should swap demes
	// per generation. .5 would have a swap
	// every other generation, on average,
	// 5 would have 5 swaps per generation,
	// etc.
	// Individuals cannot leave a deme without
	// another replacing it, in order to maintain
	// each deme's size does not change.
	// (This assertion is subject to change with
	// further analysis)
	MigrationChance float64
}

func (dg *DemeGroup) NextGeneration() bool {
	migrators := 0.0
	if dg.MigrationChance >= 1.0 {
		migrators = math.Floor(dg.MigrationChance)
	}
	if rand.Float64() < dg.MigrationChance-migrators {
		migrators += 1.0
	}

	// Shuffle the deme order
	if migrators >= 1.0 {
		for i := range dg.Demes {
			j := rand.Intn(i + 1)
			dg.Demes[i], dg.Demes[j] = dg.Demes[j], dg.Demes[i]
		}
		for i := 0; i < int(math.Floor(migrators)); i++ {
			// We pick migrators from demes in order
			// (or, randomly, as we just shuffled them)
			// and need to loop around once we've hit
			// the total deme count.
			j := i % len(dg.Demes)
			// We always migrate to the deme
			// immediately following the deme at j.
			k := (j + 1) % len(dg.Demes)
			d1 := dg.Demes[j].Members
			d2 := dg.Demes[k].Members

			// We assume all elements in the length
			// of d1 and d2 are populated with
			// individuals. This should be the case,
			// as they should be full outside of next
			// generation calls.
			m := rand.Intn(len(d1))
			n := rand.Intn(len(d2))

			// Perform the actual migration swap
			d1[m], d2[n] = d2[n], d1[m]
		}
	}
	stopEarly := false
	for i := range dg.Demes {
		stopEarly = dg.Demes[i].NextGeneration() || stopEarly
	}
	return stopEarly
}
