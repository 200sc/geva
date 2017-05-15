package eda

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

type BOA struct {
	Base
	P *pop.Population
}

func (boa *BOA) Adjust() Model {
	selected := boa.SelectLearning(boa.P)
	// Construct a network based on selected
	envs := MemberEnvs(selected)
	bn := NewBayesNet(envs)
	// Sample children from the network
	samples := bn.Sample(envs, int(boa.learningRate*float64(boa.samples)))
	// Replace worst fitnesses in boa.P with children
	boa.ReplaceLowFitnesses(boa.P, samples)
	return boa
}

func BOAModel(opts ...Option) (Model, error) {
	var err error
	boa := new(BOA)
	boa.Base, err = DefaultBase(opts...)
	boa.P = boa.Pop()
	return boa, err
}

// This bayes net construction code is based on http://www.cleveralgorithms.com/nature-inspired/probabilistic/boa.html
type BayesNet struct {
	children [][]int
	parents  [][]int
}

func NewBayesNet(samples []*env.F) *BayesNet {
	bn := new(BayesNet)
	bn.children = make([][]int, len(*samples[0]))
	bn.parents = make([][]int, len(*samples[0]))
	edgeCt := len(samples) * 3
	for i := 0; i < edgeCt; i++ {
		bestValue := 0.0
		bestFrom := -1
		bestTo := -1
		for j := range bn.children {
			value, to := bn.TryConnect(samples, j)
			if value > bestValue {
				bestValue = value
				bestFrom = j
				bestTo = to
			}
		}
		if bestValue <= 0.0 {
			break
		}
		bn.parents[bestFrom] = append(bn.parents[bestFrom], bestTo)
		bn.parents[bestTo] = append(bn.parents[bestTo], bestFrom)
	}
	return bn
}

// This is some tricky garbage the original code I based this off of
// does that makes the code difficult to understand.

func countEdges(samples []*env.F, indices []int) []int {
	counts := make([]int, int(math.Pow(2, float64(len(indices)))))
	revIndices := make([]int, len(indices))
	for i, v := range indices {
		revIndices[len(indices)-(1+i)] = v
	}
	for _, s := range samples {
		j := 0.0
		for i, v := range revIndices {
			if s.Get(v) == 1.0 {
				j += math.Pow(2, float64(i))
			}
		}
		counts[int(j)]++
	}
	return counts
}

// Because this is based on the above tricky garbage, it shared
// some of the unfortunate difficulty

func k2(i int, candidates []int, samples []*env.F) float64 {
	edgeCounts := countEdges(samples, append(candidates, i))
	total := 1.0
	for j := 0; j < len(edgeCounts)/2; j++ {
		a1, a2 := edgeCounts[j*2], edgeCounts[(j*2)+1]
		total *= (1.0 / (float64(factorial(a1+a2) + 1))) * float64(factorial(a1)) * float64(factorial(a2))
	}
	return total
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * factorial(n-1)
}

func (bn *BayesNet) TryConnect(samples []*env.F, index int) (float64, int) {
	parents := bn.ViableParents(index)
	to := -1
	bestValue := -1.0
	for _, p := range parents {
		// this should be variable
		if len(bn.children[p]) < 2 {
			v := k2(index, append(bn.children[index], p), samples)
			if v > bestValue {
				to = p
				bestValue = v
			}
		}
	}
	return bestValue, to
}

func (bn *BayesNet) ViableParents(index int) []int {
	// viable parents are those where there is no current existing path
	// between the two nodes in the net already.
	parents := []int{}
	for i := 0; i < len(bn.parents); i++ {
		if !bn.connected(index, i) {
			parents = append(parents, i)
		}
	}
	return parents
}

func (bn *BayesNet) connected(i, j int) bool {
	visited := make(map[int]bool)
	stack := []int{i}
	for len(stack) > 0 {
		next := stack[0]
		stack = stack[1:]
		visited[next] = true
		for _, c := range bn.children[next] {
			if c == j {
				return true
			}
			if _, ok := visited[c]; !ok {
				stack = append(stack, c)
			}
		}

	}
	return false
}

func (bn *BayesNet) Prob(i int, curSample *env.F, samples []*env.F) float64 {
	if len(bn.children[i]) == 0 {
		return UnivariateFromSamples(samples, i)
	}
	nodeList := []int{i}
	nodeList = append(nodeList, bn.children[i]...)
	edgeCounts := countEdges(samples, nodeList)
	j := 0.0
	revChildren := make([]int, len(bn.children[i]))
	for k := range revChildren {
		revChildren[k] = bn.children[i][len(revChildren)-(1+k)]
	}
	for k, v := range revChildren {
		if curSample.Get(v) == 1.0 {
			j += math.Pow(2, float64(k))
		}
	}
	l := j + math.Pow(2, float64(len(bn.children[i])))
	a1, a2 := float64(edgeCounts[int(l)]), float64(edgeCounts[int(j)])
	return a1 / (a1 + a2)
}

func (bn *BayesNet) SampleOrdered(samples []*env.F, ordered []int) *env.F {
	sample := env.NewF(len(*samples[0]), 0.0)
	for i := range bn.parents {
		if bn.Prob(i, sample, samples) > rand.Float64() {
			sample.Set(i, 1.0)
		} else {
			sample.Set(i, 0.0)
		}
	}
	return sample
}

func (bn *BayesNet) Topographical() []int {
	out := make([]int, len(bn.parents))
	roots := []int{}
	toSee := make([]int, len(bn.parents))
	for i, lst := range bn.parents {
		toSee[i] = len(lst)
		if len(lst) == 0 {
			roots = append(roots, i)
		}
	}
	stack := roots
	for i := range out {
		next := stack[0]
		stack := stack[1:]
		for _, c := range bn.children[next] {
			toSee[c]--
			if toSee[c] <= 0 {
				stack = append(stack, c)
			}
		}
		out[i] = next
	}
	return out

}

func (bn *BayesNet) Sample(curSamples []*env.F, n int) []*env.F {
	samples := make([]*env.F, n)
	ordered := bn.Topographical()
	for i := 0; i < n; i++ {
		samples[i] = bn.SampleOrdered(curSamples, ordered)
	}
	return samples
}
