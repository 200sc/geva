package gp

import (
	"math/rand"
)

func (gp *GP) ShrinkMutate() {

	n := gp.first

	// We have some special cases for gp.first
	// We don't want to replace gp.first with
	// a terminal action. If it already is a
	// terminal action (it shouldn't be) we
	// have to abort.
	if len(n.args) == 0 {
		return
	}

	// Avoid replacing gp.first with a terminal
	tries := 0
	for {
		i := rand.Intn(len(n.args))
		if len(n.args[i].args) != 0 {
			n = n.args[i]
			break
		}
		tries++
		if tries == len(n.args)*2 {
			return
		}
	}

	// The important part of the function
	for {

		// Continually traverse down
		// to a random child
		i := rand.Intn(len(n.args))

		// If the child has no children,
		// it is now our shrink target
		if len(n.args[i].args) == 0 {
			// shrink
			n.eval = n.args[i].eval
			n.args = n.args[i].args
			break
		} else {
			n = n.args[i]
		}
	}
}

func (gp *GP) SwapMutate() {
	nodes := gp.first.GetAllNodes()
	i := rand.Intn(len(nodes))
	children := len(nodes[i].args)
	nodes[i].eval = actions[children][rand.Intn(len(actions[children]))]
}

func (n *Node) GetAllNodes() []*Node {
	nodes := []*Node{n}
	for i := 0; i < len(n.args); i++ {
		nodes = append(nodes, n.args[i].GetAllNodes()...)
	}
	return nodes
}
