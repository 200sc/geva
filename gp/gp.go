package gp

// The principal Individual implementation for the gp package
type GP struct {
	first Node
	env   *Environment
	nodes int
}

var (
	gpOptions   GPOptions
	crossover   GPCrossover
	environment Environment
	actions     Actions
	fitness     func(n Node, inputs, outputs [][]float64) int
)

func Init(genOpt GPOptions, env Environment, cross GPCrossover,
	act Actions, f func(n Node, inputs, outputs [][]float64) int) {

	environment = env
	actions = act
	fitness = f
	gpOptions = genOpt
	crossover = cross
}

func GenerateGP(genOpt GPOptions) *GP {

	// Eventually we'll do something
	// with creation types here

	gp := new(GP)
	gp.env = &environment
	a, children := getNonZeroAction()
	gp.first = Node{
		make([]Node, children),
		a,
		gp,
	}
	gp.nodes = gp.first.GenerateTree(genOpt.MaxStartDepth, genOpt.MaxNodeCount)
	return gp
}

func (gp *GP) Print() {
	// We aren't going to print the environment,
	// because GPs are usually going to be printed
	// in large sets, and because environments are
	// usually going to be constant between GPs.
	gp.first.Print("", true)
}

func (gp *GP) CanCrossover(other population.Individual) bool {
	switch other.(type) {
	default:
		return false
	case *GP:
		return true
	}
}

func (gp *GP) Crossover(other population.Individual) population.Individual {
	// I'm relatively confident the most obvious way to crossover
	// gp trees is to use point crossover.
	gp2 := other.(*GP)
	return crossover.Crossover(gp, gp2)
}

func (gp *GP) Fitness(input, expected [][]float64) int {
	return fitness(gp.first, input, expected)
}

func (gp *GP) Mutate() {
	// We refer to our mutation chances here, and run
	// through the tree with those chances. We probably
	// want to start at the leaves.
}
