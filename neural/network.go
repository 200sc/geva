package neural

// Store functions exclusive to networks.

import (
	"fmt"
	"goevo/population"
	"math"
	"math/rand"
)

func (nn *Network) Crossover(other population.Individual) population.Individual {
	// Assert that other is a Network
	nn2 := other.(*Network)
	// Perform a crossover method as defined by global settings
	return crossover.Crossover(nn, nn2)
}

func (nn *Network) CanCrossover(other population.Individual) bool {
	return true
}

/**
 * Copy a network
 */
func (nn *Network) Copy() Network {
	newBody := nn.Body.Copy()
	return Network{
		Body:      newBody,
		Activator: nn.Activator,
	}
}

/**
 * Copy a network body
 */
func (nn_p *Body) Copy() Body {

	nn := *nn_p

	var newNetwork Body

	for i := range nn {
		newNetwork = append(newNetwork, make([]Neuron, len(nn[i])))
		copy(newNetwork[i], nn[i])
	}

	return newNetwork
}

func (genOpt NetworkGenerationOptions) Generate() *Network {
	return GenerateNetwork(&genOpt)
}

/**
 * Convert generation options into
 * a new neural network
 */
func GenerateNetwork(nOpt_p *NetworkGenerationOptions) *Network {

	nnOpt := *nOpt_p
	cOpt := *nnOpt.ColumnOptions

	nn := make(Body, 0)

	// Set up the input column
	inputColumn := make([]Neuron, nnOpt.Inputs)

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := make([]Neuron, nnOpt.Outputs)

	nn = append(nn, outputColumn)

	// reset the input column to give it axons
	for i := 0; i < len(inputColumn); i++ {
		nn[0][i] = make(Neuron, len(outputColumn))
		nn.replaceNeuron(0, i, cOpt.DefaultAxonWeight)
	}

	columnCount := rand.Intn(nnOpt.MaxColumns-nnOpt.MinColumns) + nnOpt.MinColumns

	for i := 0; i < columnCount; i++ {
		nn.addColumn(&cOpt)
	}

	for i := 0; i < nnOpt.BaseMutations; i++ {
		nn.Mutate(&nnOpt.NetworkMutationOptions)
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
func (modNet_p *Network) Run(Inputs []float64) []float64 {

	modNet := *modNet_p
	act := modNet.Activator
	nn := modNet.Body

	doneCh := make(chan networkOutput)

	channels := make([][]chan float64, len(nn))
	for x, col := range nn {
		channels = append(channels, []chan float64{})
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
					// Inputs for our value.
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
		ch <- Inputs[i]
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
func (n Network) Fitness(inputs, expected [][]float64) int {
	fitness := 1.0
	for i := range inputs {
		output := n.Run(inputs[i])
		for j := range output {
			fitness += math.Abs(output[j] - expected[i][j])
		}
	}
	return int(math.Ceil(fitness))
}
