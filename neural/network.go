package neural

import (
	"fmt"
	"github.com/200sc/geva/pop"
	"math/rand"
)

var (
	ngo       NetworkGenerationOptions
	crossover NeuralCrossover
	fitness   FitnessFunc
)

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

func (nn *Network) Crossover(other pop.Individual) pop.Individual {
	return crossover.Crossover(nn, other.(*Network))
}

func (nn *Network) CanCrossover(other pop.Individual) bool {
	switch other.(type) {
	case *Network:
		return true
	default:
		return false
	}
}

func (nn *Network) Copy() *Network {
	return &Network{
		Body:      nn.Body.Copy(),
		Activator: nn.Activator,
	}
}

func (nn_p *Body) Copy() Body {

	nn := *nn_p
	var newNetwork Body

	for i := range nn {
		newNetwork = append(newNetwork, make([]Neuron, len(nn[i])))
		copy(newNetwork[i], nn[i])
	}

	return newNetwork
}

func GeneratePopulation(opt interface{}, popSize int) []pop.Individual {
	nnOpt := opt.(NetworkGenerationOptions)
	members := make([]pop.Individual, popSize)
	for j := 0; j < popSize; j++ {
		members[j] = GenerateNetwork(nnOpt)
	}
	return members
}

/**
 * Convert generation options into
 * a new neural network
 */
func GenerateNetwork(nnOpt NetworkGenerationOptions) *Network {

	cOpt := nnOpt.ColumnOptions

	nn := make(Body, 0)

	// Set up the input column
	inputColumn := make([]Neuron, nnOpt.MaxInputs)

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := make([]Neuron, nnOpt.MaxOutputs)

	nn = append(nn, outputColumn)

	// reset the input column to give it axons
	for i := 0; i < len(inputColumn); i++ {
		nn[0][i] = make(Neuron, len(outputColumn))
		nn.replaceNeuron(0, i, cOpt.DefaultAxonWeight)
	}

	columnCount := rand.Intn(nnOpt.MaxColumns-nnOpt.MinColumns) + nnOpt.MinColumns

	for i := 0; i < columnCount; i++ {
		nn.addColumn(cOpt)
	}

	for i := 0; i < nnOpt.BaseMutations; i++ {
		nn.Mutate(nnOpt.NetworkMutationOptions)
	}

	return &Network{
		Body:      nn,
		Activator: MutateActivator(nnOpt.ActivatorOptions),
	}
}

/**
 * Print a network
 */
func (nn Network) Print() {
	for _, col := range nn.Body {
		for _, n := range col {
			fmt.Print(n.String())
		}
		fmt.Println("")
	}
	GraphPrintActivator(nn.Activator)
	fmt.Println("")
}

/**
 * Run some input through a neural network.
 * This returns the network's output column.
 */
func (modNet_p *Network) Run(inputs []float64) []float64 {

	modNet := *modNet_p
	act := modNet.Activator
	nn := modNet.Body

	doneCh := make(chan networkOutput)

	channels := make([][]chan float64, len(nn))
	for x, col := range nn {
		channels[x] = []chan float64{}
		for range col {
			channels[x] = append(channels[x], make(chan float64))
		}
	}
	for x, col := range nn {
		for y, neuron := range col {
			var l int
			if x == 0 {
				l = 1
			} else {
				l = len(channels[x-1])
			}
			if x == len(nn)-1 {
				go func(inputChannel chan float64, doneCh chan networkOutput,
					inputLength int, y int, fn ActivatorFunc) {

					out := 0.0

					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					out = fn(out)

					doneCh <- networkOutput{out, y}

				}(channels[x][y], doneCh, l, y, act)
			} else {

				go func(n Neuron, inputChannel chan float64,
					channelColumn []chan float64, inputLength int, fn ActivatorFunc) {

					out := 0.0

					// this system has each neuron
					// know the weights of the following layer,
					// instead of knowing where they output to.
					// It also has each neuron send to every
					// neuron in the following layer-- a 0
					// weight is equivalent to not accepting
					// some input.
					//
					// At this stage, that means we can
					// just sum all of our already-weighted
					// inputs for our value.
					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					out = fn(out)

					// As above, we apply the next column's
					// weights as we send them off.
					for i, weight := range n {
						channelColumn[i] <- out * weight
					}
				}(neuron, channels[x][y], channels[x+1], l, act)
			}
		}
	}

	// Send the first row their initial values
	for i, ch := range channels[0] {
		if i >= len(inputs) {
			ch <- 0.0
		} else {
			ch <- inputs[i]
		}
	}

	// We need to wait here, on the last columns being populated
	output := make([]float64, len(nn[len(nn)-1]))
	for i := 0; i < len(nn[len(nn)-1]); i++ {
		recieved := <-doneCh
		output[recieved.index] = recieved.value
	}
	close(doneCh)

	return output
}

/**
 * Get a set of slices
 * in the same shape as
 * a network body's slices.
 */
func (b Body) CopyStructure() Body {
	body := make(Body, len(b))
	for i := 0; i < len(b); i++ {
		body[i] = make([]Neuron, len(b[i]))
	}
	return body
}

/**
 * Evaluate the fitness of a network
 * low fitness is good, high fitness is bad.
 */
func (n *Network) Fitness(inputs, expected [][]float64) int {
	return fitness(n, inputs, expected)
}
