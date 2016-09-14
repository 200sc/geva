package gp

import (
	"math/rand"
)

type Node struct {
	args []Node
	eval func(*GP, ...Node) int
	// Each node has to implicitly know what GP it is a part of
	// so it can refer to the environment of the GP if it is a
	// leaf node.
	gp *GP
}

func (n *Node) GenerateTree(depth, nodeLimit int) int {

	nodes := 1

	for i := 0; i < len(n.args); i++ {
		// Get a terminal action if...
		if depth == 0 || nodes >= nodeLimit || rand.Float64() < 0.5 {
			a, _ := getAction(0)
			n.args[i] = Node{
				make([]Node, 0),
				a,
				n.gp,
			}
			nodes++
		} else {
			a, children := getNonZeroAction()
			n.args[i] = Node{
				make([]Node, children),
				a,
				n.gp,
			}
			nodes += n.args[i].GenerateTree(depth-1, nodeLimit-nodes)
		}
	}

	return nodes
}

// This is more or less shorthand
func Eval(n Node) int {
	return n.eval(n.gp, n.args...)
}

//http://stackoverflow.com/questions/4965335/how-to-print-binary-tree-diagram
// Original credit, as per above source, to whoever wrote tree for linux
func (n *Node) Print(prefix string, isTail bool) {
	// Uh we need a way to identify what the
	// action a node has in a printable manner
	// and we can't do that yet
	s := prefix
	if isTail {
		s += "└──"
		prefix += "    "
	} else {
		s += "├──"
		prefix += "│   "
	}
	// Add identifier here
	s += GetOperatorName(n.eval)
	fmt.Println(s)
	for i := 0; i < len(n.args)-1; i++ {
		n.args[i].Print(prefix, false)
	}
	if len(n.args) > 0 {
		n.args[len(n.args)-1].Print(prefix, true)
	}
}
