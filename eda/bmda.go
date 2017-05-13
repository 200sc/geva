package eda

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// BMDA represents the Bivariate Marginal Distribution Algorithm
type BMDA struct {
	UMDA
	BF      FullBivariateEnv
	LastPop *pop.Population
}

// ChiSquared on a bmda calculates a chi^2 value
// where the observed values are the unconditional bivariate
// probabilities and the expected values are the univariate
// probabilities multipled together.
func (bmda *BMDA) ChiSquared(a, b int) float64 {

	// p(a=t)
	at := bmda.F.Get(a)
	bt := bmda.F.Get(b)
	af := 1 - at
	bf := 1 - bt

	// p(a=t|b=t)
	catbt := bmda.BF[a].bf[0].Get(b)
	catbf := bmda.BF[a].bf[1].Get(b)
	cafbt := 1 - catbt
	cafbf := 1 - catbf

	// p(a=t|b=t) = p(a=t,b=t) / p(b=t) so
	// p(a=t,b=t)
	atbt := catbt * bt
	atbf := catbf * bf
	afbt := cafbt * bt
	afbf := cafbf * bf

	chi2 := 0.0
	fsamp := float64(bmda.samples)

	xptd1 := fsamp * at * bt
	if xptd1 != 0 {
		chi2 += math.Pow((fsamp*atbt)-xptd1, 2) / xptd1
	}

	xptd2 := fsamp * at * bf
	if xptd2 != 0 {
		chi2 += math.Pow((fsamp*atbf)-xptd2, 2) / xptd2
	}

	xptd3 := fsamp * af * bt
	if xptd3 != 0 {
		chi2 += math.Pow((fsamp*afbt)-xptd3, 2) / xptd3
	}

	xptd4 := fsamp * af * bf
	if xptd4 != 0 {
		chi2 += math.Pow((fsamp*afbf)-xptd4, 2) / xptd4
	}

	//fmt.Println("Chi2", chi2, a, b)
	return chi2
}

// Adjust on a BMDA is incomplete
func (bmda *BMDA) Adjust() Model {
	// Create dependency forest
	roots := []int{}
	children := make([][]int, bmda.length)
	for i := 0; i < bmda.length; i++ {
		children[i] = []int{}
	}
	available := bmda.GenIndices()
	used := []int{}
	chi2Memo := make(map[int]map[int]float64)
	//fmt.Println("Starting forest generation")
	for {
		// choose a random index
		i := rand.Intn(len(available))
		chosen := available[i]
		// create a new tree in the forest starting at chosen
		roots = append(roots, chosen)
		for len(available) > 0 {
			// remove the chosen index from available
			available = append(available[:i], available[i+1:]...)
			used = append(used, chosen)

			// Find the most dependant pairing given our set of used
			// and available indices
			maxChi2 := 0.0
			var parent int
			// This loop is incredibly expensive
			// It is probably faster to calculate all n^2 chiSquared values
			// ahead of time
			for j, v := range available {
				for _, v2 := range used {
					if _, ok := chi2Memo[v]; !ok {
						chi2Memo[v] = make(map[int]float64)
					}
					if _, ok := chi2Memo[v][v2]; !ok {
						chi2Memo[v][v2] = bmda.ChiSquared(v, v2)
					}
					if chi2Memo[v][v2] > maxChi2 {
						maxChi2 = chi2Memo[v][v2]
						chosen = v
						i = j
						parent = v2
					}
				}
			}

			// If no dependencies exist break to the outer loop.
			if maxChi2 < 3.84 {
				break
			}
			children[parent] = append(children[parent], chosen)
			//fmt.Println("Length of available:", len(available))
		}
		//fmt.Println("Length of available:", len(available))
		if len(available) == 0 {
			break
		}
	}
	//fmt.Println("End forest")
	// fmt.Println(roots, children)
	// Generate new population from forest and frequencies
	newPop := bmda.BMDAPop(roots, children)
	// Combine previous population and new population by direct replacement
	// of the worst elements

	// Get the samples from the population
	samples := make([]*env.F, len(bmda.LastPop.Members))
	for i, mem := range bmda.LastPop.Members {
		samples[i] = mem.(*UMDAIndividual).F
	}

	//fmt.Println(bmda.LastPop.Members)
	//fmt.Println(samples)
	//time.Sleep(1 * time.Second)

	// Sort the samples by their fitnesses
	_, samples = SampleFitnesses(bmda, samples)
	// Replace the end samples (with the worst fitness) with the new population
	for i, e := range newPop {
		samples[len(samples)-(i+1)] = e
	}

	// Reassign samples to the population
	for i := range bmda.LastPop.Members {
		bmda.LastPop.Members[i] = &UMDAIndividual{samples[i]}
	}
	//fmt.Println(bmda.LastPop.Members)

	bmda.UpdateFromPop()
	for _, be := range bmda.BF {
		be.bf[0].Mutate(bmda.mutationRate, bmda.fmutator)
		be.bf[1].Mutate(bmda.mutationRate, bmda.fmutator)
	}
	return bmda
}

// BMDAModel returns an initialized BMDA
func BMDAModel(opts ...Option) (Model, error) {
	var err error
	bmda := new(BMDA)
	bmda.Base, err = DefaultBase(opts...)
	// Generate initial population
	bmda.LastPop = bmda.Pop()
	//fmt.Println(bmda.LastPop.Members)
	bmda.UpdateFromPop()
	return bmda, err
}

// How do you sample a 2D bivariate array
// No you use the foest we just spent way too long making
// So we pick random (all) roots and go down through all of their children
// using their parents as the elements they are dependant on
func (bmda *BMDA) BMDAPop(roots []int, children [][]int) []*env.F {

	// Create environments to sample from
	tenv := env.NewF(bmda.length, 0.0)
	fenv := env.NewF(bmda.length, 0.0)

	// Get the samples from the population
	samples := make([]*env.F, len(bmda.LastPop.Members))
	for i, mem := range bmda.LastPop.Members {
		samples[i] = mem.(*UMDAIndividual).F
	}

	for _, root := range roots {
		// Roots just use the univariate probability
		tenv.Set(root, UnivariateFromSamples(samples, root))
		bmda.SetChildren(samples, tenv, fenv, children, root)
	}

	return bmda.NSamples(bmda.learningSamples, tenv, fenv, roots, children)
}

func (bmda *BMDA) NSamples(n int, tenv, fenv *env.F, roots []int, children [][]int) []*env.F {
	samples := make([]*env.F, n)
	for i := 0; i < n; i++ {
		samples[i] = bmda.GetSample(tenv, fenv, roots, children)
	}
	return samples
}

func (bmda *BMDA) GetSample(tenv, fenv *env.F, roots []int, children [][]int) *env.F {
	senv := env.NewF(bmda.length, 0.0)
	for _, root := range roots {
		senv.Set(root, tenv.GetBinary(root))
		bmda.SetChildSample(senv, tenv, fenv, children, root)
	}
	return senv
}
func (bmda *BMDA) SetChildSample(senv, tenv, fenv *env.F, children [][]int, parent int) {
	for _, child := range children[parent] {
		e := ConditionedBSEnv(senv, tenv, fenv, parent)
		senv.Set(child, e.GetBinary(child))
		bmda.SetChildSample(senv, tenv, fenv, children, child)
	}
}
func (bmda *BMDA) SetChildren(samples []*env.F, tenv, fenv *env.F, children [][]int, parent int) {
	// If parent has no children, this is a NOP
	for _, child := range children[parent] {
		ptt, ptf := BitStringBivariate(samples, child, parent)
		tenv.Set(child, ptt)
		fenv.Set(child, ptf)
		bmda.SetChildren(samples, tenv, fenv, children, child)
	}
}

func (bmda *BMDA) UpdateFromPop() {
	bmda.LastPop.Size = bmda.learningSamples
	// Select parents?
	// I do not know why these are called parents in the pseudocode for the
	// BMDA paper, they're just the better members of the population
	subPop := bmda.selection.Select(bmda.LastPop)
	envs := make([]*env.F, len(subPop))
	for i, mem := range subPop {
		envs[i] = mem.(*UMDAIndividual).F
	}
	// Calculate univariate and bivariate frequencies
	// Univariate
	bmda.F.SetAll(0.0)
	bmda.F.AddF(envs...)
	bmda.F.Divide(float64(len(subPop)))
	// At this point bmda.F holds the univariate frequencies
	// BF is the conditional bivariate frequencies
	bmda.BF = NewFullBSBivariateEnv(envs)
	bmda.LastPop.Size = bmda.samples
}
