package gp

import (
	"fmt"
	"math"
	"math/rand"
)

type Node struct {
	args []*Node
	eval Action
	// Each node has to implicitly know what GP it is a part of
	// so it can refer to the environment of the GP if it is a
	// leaf node.
	gp *GP
}

func NewNode(a Action, children int, gp *GP) *Node {
	return &Node{
		make([]*Node, children),
		a,
		gp,
	}
}

func (n *Node) GenerateTree(depth, nodeLimit int) int {

	nodes := 1

	for i := 0; i < len(n.args); i++ {
		// Get a terminal action if...
		if depth == 0 || nodes >= nodeLimit || rand.Float64() < 0.5 {
			a := getZeroAction()
			n.args[i] = NewNode(a, 0, n.gp)
			nodes++
		} else {
			a, children := getNonZeroAction()
			n.args[i] = NewNode(a, children, n.gp)
			nodes += n.args[i].GenerateTree(depth-1, nodeLimit-nodes)
		}
	}

	return nodes
}

// This is more or less shorthand
func Eval(n *Node) int {
	return n.eval.op(n.gp, n.args...)
}

//http://stackoverflow.com/questions/4965335/how-to-print-binary-tree-diagram
// Original credit, as per above source, to whoever wrote tree for linux
func (n *Node) Print(prefix string, isTail bool) {

	s := prefix
	if isTail {
		s += "└──"
		prefix += "    "
	} else {
		s += "├──"
		prefix += "│   "
	}
	// Add identifier here
	s += n.eval.name
	fmt.Println(s)
	for i := 0; i < len(n.args)-1; i++ {
		n.args[i].Print(prefix, false)
	}
	if len(n.args) > 0 {
		n.args[len(n.args)-1].Print(prefix, true)
	}
}

func (n *Node) GetRandomNode() (*Node, *Node) {

	// We assume the depth of the GP (which is not stored)
	// is about log2() of the nodes in the GP.
	// using that we try to have a 50% chance to take a node
	// at around halfway down the tree.

	apxDepth := math.Floor(math.Log2(float64(n.gp.nodes)))

	curDepth := 1

	// Don't choose this node. Choose from first tier children, at least
	prevN := n
	n = n.RandChild()

	for {
		if len(n.args) == 0 {
			return n, prevN
		}
		chance := float64(curDepth) / float64(apxDepth)

		if rand.Float64() < chance {
			return n, prevN
		}

		curDepth++
		prevN = n
		n = n.RandChild()
	}
}

func (n *Node) RandChild() *Node {
	return n.args[rand.Intn(len(n.args))]
}

func (n *Node) Copy(gp *GP) *Node {
	newNode := new(Node)
	newNode.eval = n.eval
	newNode.gp = gp
	newNode.args = make([]*Node, len(n.args))
	for i, child := range n.args {
		newNode.args[i] = child.Copy(gp)
	}
	return newNode
}

func (n *Node) Size() int {
	size := 1
	for _, child := range n.args {
		size += child.Size()
	}
	return size
}
