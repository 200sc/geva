package gp

type GPOptions struct {
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

func (gpo *GPOptions) MaxNodeCount(mnc int) *GPOptions {
	gpo.MaxNodeCount = nc
	return gpo
}

func (gpo *GPOptions) MaxStartDepth(msd int) *GPOptions {
	gpo.MaxStartDepth = msd
	return gpo
}

func (gpo *GPOptions) MaxDepth(md int) *GPOptions {
	gpo.MaxDepth = md
	return gpo
}

func (gpo *GPOptions) SwapMutationChance(smc float64) *GPOptions {
	gpo.SwapMutationChance = smc
	return gpo
}

func (gpo *GPOptions) ShrinkMutationChance(smc float64) *GPOptions {
	gpo.ShrinkMutationChance = smc
	return gpo
}
