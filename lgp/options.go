package lgp

type LGPOptions struct {
	MaxActionCount  int
	MaxStartActions int
	MinStartActions int

	SwapMutationChance   float64
	ValueMutationChance  float64
	ShrinkMutationChance float64
	ExpandMutationChance float64
	MemMutationChance    float64
}
