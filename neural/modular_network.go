package neural

import (
	"fmt"
	"math"
	"math/rand"
)

type ModularNetworkOutput struct {
	value float64
	index int
}

// An activator function just maps float values to other
// float values. The function can be as simplistic or complicated
// as desired-- eventually a set of common activators will be
// collected.
type ActivatorFunc func(float64) float64

type ModularBody [][]ModularNeuron

type ModularNetwork struct {
	Activator ActivatorFunc
	Body      ModularBody
}

func (modNet_p *ModularNetwork) Copy() ModularNetwork {
	newBody := modNet_p.Body.Copy()
	return ModularNetwork{
		Body:      newBody,
		Activator: modNet_p.Activator,
	}
}

/**
 * Take a network and duplicate it
 */
func (nn_p *ModularBody) Copy() ModularBody {

	nn := *nn_p

	var newNetwork ModularBody

	for i := range nn {
		newNetwork = append(newNetwork, make([]ModularNeuron, len(nn[i])))
		copy(newNetwork[i], nn[i])
	}

	return newNetwork
}

func GenerateModularNetwork(nOpt_p *ModularNetworkGenerationOptions) *ModularNetwork {

	nnOpt := *nOpt_p
	cOpt := *nnOpt.ColumnOptions

	nn := make(ModularBody, 0)

	// Set up the input column
	inputColumn := make([]ModularNeuron, nnOpt.Inputs)

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := make([]ModularNeuron, nnOpt.Outputs)

	nn = append(nn, outputColumn)

	// reset the input column to give it axons
	for i := 0; i < len(inputColumn); i++ {
		nn[0][i] = make(ModularNeuron, len(outputColumn))
		nn.replaceNeuron(0, i, cOpt.DefaultAxonWeight)
	}

	columnCount := rand.Intn(nnOpt.MaxColumns-nnOpt.MinColumns) + nnOpt.MinColumns

	for i := 0; i < columnCount; i++ {
		nn = *(nn.addColumn(&cOpt))
	}

	for i := 0; i < nnOpt.BaseMutations; i++ {
		nn = *(nn.Mutate(&nnOpt.ModularNetworkMutationOptions))
	}

	modNet := ModularNetwork{
		Body:      nn,
		Activator: nnOpt.Activator,
	}

	return &modNet
}

// Todo: Print Activator
func (nn ModularNetwork) Print() {
	for _, col := range nn.Body {
		for _, n := range col {
			fmt.Print(n.String())
		}
		fmt.Println("")
	}
	fmt.Println("")
}

/**
 * Run some input through a neural network.
 * This returns the network's output column.
 */
func (modNet_p *ModularNetwork) Run(Inputs []float64) []float64 {

	modNet := *modNet_p
	act := modNet.Activator
	nn := modNet.Body

	doneCh := make(chan ModularNetworkOutput)

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
				go func(inputChannel chan float64, doneCh chan ModularNetworkOutput,
					inputLength int, y int, fn ActivatorFunc) {

					out := 0.0

					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					out = fn(out)

					doneCh <- ModularNetworkOutput{out, y}

				}(channels[x][y], doneCh, l, y, act)
			} else {

				go func(n ModularNeuron, inputChannel chan float64,
					channelColumn []chan float64, inputLength int, fn ActivatorFunc) {

					out := 0.0

					// Compared to network.go,
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

func (b ModularBody) CopyStructure() ModularBody {
	body := make(ModularBody, len(b))
	for i := 0; i < len(b); i++ {
		body[i] = make([]ModularNeuron, len(b[i]))
	}
	return body
}

// High fitness is bad, and vice versa.
func (n ModularNetwork) Fitness(Inputs, expected [][]float64) int {
	fitness := 1.0
	for i := range Inputs {
		output := n.Run(Inputs[i])
		for j := range output {
			fitness += math.Abs(output[j] - expected[i][j])
		}
	}
	return int(math.Ceil(fitness))
}
