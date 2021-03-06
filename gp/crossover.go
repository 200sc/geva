package gp

import (
	"math/rand"
)

type GPCrossover interface {
	Crossover(a, b *GP) *GP
}

type PointCrossover struct{}

// Neural nets use a sort of 'number of points'
// to combine two networks at. We could do that here,
// but we wouldn't have a way, in this tree structure,
// to avoid overlapping our crossovers. If we were picking
// from suitably high up to take trees from we'd almost
// always have duplicate structures in our new tree from
// the result of more than one crossover point being
// used. Even if that didn't happen, the crossover would
// be equivalent to performing PointCrossover twice or
// however many times. There is something to be said for
// having a low chance of a very significant mutation,
// though.

func (pc PointCrossover) Crossover(a, b *GP) *GP {

	g1 := a.First
	g2 := b.First

	c := new(GP)

	g3 := g1.Copy(c)

	// Find a random point in both networks.
	// Replace what exists at g1point with g2point.
	node1, parent := g3.GetRandomNode()
	node2, _ := g2.GetRandomNode()

	// If the node chosen has no arguments,
	// we put this branch somewhere in the
	// node's parents.
	//
	// This assumes that GPs are never just one
	// no-argument instruction.
	if len(node1.args) > 0 {
		i := rand.Intn(len(node1.args))
		node1.args[i] = node2.Copy(c)
	} else {
		i := rand.Intn(len(parent.args))
		parent.args[i] = node2.Copy(c)
	}

	c.Env = a.Env
	c.Mem = a.Mem.Copy()
	c.First = g3
	c.Nodes = c.First.Size()
	for c.Nodes > gpOptions.MaxNodeCount {
		c.ShrinkMutate()
		c.Nodes--
	}
	return c
}
