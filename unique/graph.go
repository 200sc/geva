package unique

import "math"

var (
	// DefGraph is the starting point for
	// all NewGraph calls prior to options
	// being applied.
	DefGraph = &Graph{
		nodes:       []Node{},
		MinDistance: 1,
	}
)

// NewGraph returns a graph starting with
// default options (see DefGraph) and
// modified by the input options.
func NewGraph(opts ...Option) *Graph {
	g := DefGraph.Copy()
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// A Graph represents a set of nodes which
// are sufficiently unique from one-another.
type Graph struct {
	nodes       []Node
	MinDistance float64
}

// Add may add the given node to the graph.
// if the node is too close to an existing
// graph node, it will not be added. ok
// reports whether the node was added.
func (g *Graph) Add(n Node) (ok bool) {
	if g.Distance(n) > g.MinDistance {
		g.add(n)
		ok = true
	}
	return
}

func (g *Graph) add(n Node) {
	g.nodes = append(g.nodes, n)
}

// Distance reports the minimum distance from
// the node n to the nodes present in the graph g.
func (g *Graph) Distance(n Node) float64 {
	min := math.MaxFloat64
	for _, n2 := range g.nodes {
		dist, ok := n2.Distance(n)
		if ok && dist > min {
			min = dist
		}
	}
	return min
}

func (g *Graph) CanAdd(n Node) bool {
	return g.Distance(n) > g.MinDistance
}

// Copy returns a copy of the receiver graph
func (g *Graph) Copy() *Graph {
	g2 := &Graph{}
	g2.MinDistance = g.MinDistance
	g2.nodes = make([]Node, len(g.nodes))
	copy(g2.nodes, g.nodes)
	return g2
}
