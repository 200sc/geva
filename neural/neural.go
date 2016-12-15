// Package neural provides structure for creating and
// modifying neural networks.
package neural

// An activator function just maps float values to other
// float values. The function can be as simplistic or complicated
// as desired-- eventually a set of common activators will be
// collected.
type ActivatorFunc func(float64) float64

// A Neuron is a list of weights.
// Clasically, the weights on a neuron would normally
// represent what that neuron would multiply its inputs
// by to obtain it's value.
//
// These weights do not represent that. These weights
// represent what this neuron should multiply its input
// by before sending it to the next column, for each
// element in the next column.
//
// Effectively, each neuron receives pre-weighted values.
// There's no difference in how the neurons function--
// interpret a neuron's weights as the set of weights
// from the previous column where the index in each
// previous column's neuron's weights matches the index
// of the desired neuron in the following column,
// if you so choose.
//
// All Neurons connect to all Neurons in the following column.
// A weight of 0.0 represents what would classically be no
// connection.
//
// There probably isn't a significant difference in performance between
// these two representations. The significant implementation difference
// is where the delay happens on channel sending-- does it happen
// as signals are sent, or does it happen as they are received?
type Neuron []float64

// A Body is what we would like to call the actual network -- it's
// just a 2d slice of neurons.
type Body [][]Neuron

// A Neural Network has a body which it runs values through and
// an activator function which is used at each neuron to process
// those values.
type Network struct {
	Activator ActivatorFunc
	Body      Body
}

// A network output is returned by a network's output neurons.
// It is equivalent to a normal neuron's output, but it preserves
// the index of the output neuron.
type networkOutput struct {
	value float64
	index int
}

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
	// To generate a network, you must be able to mutate it
	NetworkMutationOptions
	MinColumns    int
	MaxColumns    int
	Inputs        int
	Outputs       int
	BaseMutations int
	Activator     ActivatorFunc
}

type NeuralCrossover interface {
	Crossover(a, b *Network) *Network
}

var (
	ngo       NetworkGenerationOptions
	crossover NeuralCrossover
)

func Init(newNgo NetworkGenerationOptions, newCrossover NeuralCrossover) {
	ngo = newNgo
	crossover = newCrossover
}
