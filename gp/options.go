package gp

type Options struct {
	MaxNodeCount  int
	MaxStartDepth int
	MaxDepth      int
	// These mutation chances are rolled once per GP
	//
	// Swap mutation is a swap of one node's action
	// to another action from the pool with the same
	// number of arguments.
	// (This is based off of GPJPP. Question: why is
	// there no mutation swap that can change the number
	// of arguments?)
	SwapMutationChance float64
	// Shrink Mutation picks some leaf node and replaces
	// its parent with it.
	ShrinkMutationChance float64
}
