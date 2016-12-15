// Package neural provides structure for creating and
// modifying neural networks.
package neural

// Mutation and generation options follow.

type FloatMutationOptions struct {
	MutChance    float64
	MutMagnitude float64
	MutRange     int
	// Zeroing out a neuron is useful for making it easier
	// to remove an unnecessary connection in the network.
	ZeroOutChance float64
}

type ColumnGenerationOptions struct {
	MinSize           int
	MaxSize           int
	DefaultAxonWeight float64
}

type ActivatorMutationOptions []ActivatorFunc

type NetworkMutationOptions struct {
	WeightOptions    FloatMutationOptions
	ColumnOptions    ColumnGenerationOptions
	ActivatorOptions ActivatorMutationOptions
	// checked per column
	NeuronReplacementChance float64
	NeuronAdditionChance    float64
	WeightSwapChance        float64
	// checked per network
	ColumnRemovalChance     float64
	ColumnAdditionChance    float64
	NeuronMutationChance    float64
	ActivatorMutationChance float64
}

type NetworkGenerationOptions struct {
	NetworkMutationOptions
	MinColumns    int
	MaxColumns    int
	MaxInputs     int
	MaxOutputs    int
	BaseMutations int
}
