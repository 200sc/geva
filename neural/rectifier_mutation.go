package neural

import (
	"math/rand"
)

type RectifierNetworkMutationOptions struct {
	WeightOptions *FloatMutationOptions
	ColumnOptions *RectifierColumnGenerationOptions
	// per column
	NeuronReplacementChance float64
	NeuronAdditionChance    float64
	WeightSwapChance        float64
	// per network
	ColumnRemovalChance  float64
	ColumnAdditionChance float64
	NeuronMutationChance float64
}
type RectifierNetworkGenerationOptions struct {
	RectifierNetworkMutationOptions
	MinColumns    int
	MaxColumns    int
	Inputs        int
	Outputs       int
	BaseMutations int
}
type RectifierColumnGenerationOptions struct {
	MinSize           int
	MaxSize           int
	DefaultAxonWeight float64
}

/**
 * Mutate this neuron.
 */
func (n *RectifierNeuron) mutate(wOpt_p *FloatMutationOptions) RectifierNeuron {

	newNeuron := make(RectifierNeuron, len(*n))
	for i, weight := range *n {
		newNeuron[i] = mutateFloat(weight, *wOpt_p)
	}

	return newNeuron
}

/**
 * Mutate this network.
 */
func (nn *RectifierNetwork) Mutate(mOpt_p *RectifierNetworkMutationOptions) *RectifierNetwork {

	mOpt := *mOpt_p

	newNetwork := nn.Copy()

	if rand.Float64() < mOpt.ColumnRemovalChance {
		// We currently only remove the len-1th column
		// We can't remove a column if we have just one
		// hidden layer or the output column to remove
		if len(newNetwork) > 2 {
			newNetwork = *(newNetwork.removeColumn(mOpt.ColumnOptions))
		}
	}

	if rand.Float64() < mOpt.ColumnAdditionChance {
		// We currently only add a column in the len-1th space
		// In the future a check against some maximum column count
		// should be here
		newNetwork = *(newNetwork.addColumn(mOpt.ColumnOptions))
	}

	for i := 0; i < len(newNetwork)-1; i++ {

		col := newNetwork[i]

		// Swap two axons connecting this column
		// to the next column
		if rand.Float64() < mOpt.WeightSwapChance {

			// We can't make a meaningful swap
			// if the column only has one neuron
			if len(col) > 1 {
				// Get our neuron index
				neuronIndex := rand.Intn(len(col))

				// Get our first axon index
				axonIndex1 := rand.Intn(len(newNetwork[i+1]))

				// Get our second axon index
				axonIndex2 := rand.Intn(len(newNetwork[i+1]))
				// If our second random neuron is same
				// as our first, we interpret that as
				// taking the neuron just following our first
				// as our second
				if axonIndex2 == axonIndex1 {
					axonIndex2 = (axonIndex2 + 1) % len(newNetwork[i+1])
				}

				newNetwork.swapWeights(i, neuronIndex, axonIndex1, axonIndex2)
			}
		}

		if i != 0 {
			if rand.Float64() < mOpt.NeuronAdditionChance {
				// In the future a check on some max length
				// for a column should exist.
				newNetwork.addNeuron(i, (*mOpt.ColumnOptions).DefaultAxonWeight)
			}

			if rand.Float64() < mOpt.NeuronReplacementChance {
				neuronIndex := rand.Intn(len(col))
				newNetwork.replaceNeuron(i, neuronIndex, (*mOpt.ColumnOptions).DefaultAxonWeight)
			}
		}
	}

	// Mutate individual neurons
	for x, column := range newNetwork {
		for y, neuron := range column {

			// For performance in the future,
			// We should replace having x * y random calls
			// with picking an average number of neurons
			// to always mutate where the average is
			// calculated by x * y * MutChance. Then
			// we'd make x * y * MutChance random index
			// to-mutate determinations.
			//
			// This also applies to rand.Float64() calls
			// above.
			if rand.Float64() < mOpt.NeuronMutationChance {
				newNetwork[x][y] = neuron.mutate(mOpt.WeightOptions)
			}
		}
	}

	return &newNetwork
}

/**
 * Removes a neuron from an index and places a new neuron there in its place.
 * Effectively resetNeuron + addNeuron, if addNeuron took an index.
 */
func (nn_p *RectifierNetwork) replaceNeuron(columnIndex, neuronIndex int, DefaultAxonWeight float64) {

	nn := *nn_p

	neuron := nn[columnIndex][neuronIndex]

	for i := 0; i < len(nn[columnIndex+1]); i++ {
		neuron[i] = DefaultAxonWeight
	}
}

/**
 * Add a neuron to the end of a column.
 */
func (nn_p *RectifierNetwork) addNeuron(columnIndex int, DefaultAxonWeight float64) {

	nn := *nn_p

	newNeuron := make(RectifierNeuron, len(nn[columnIndex+1]))

	// Set this new neuron's weights for the next column
	for i := 0; i < len(nn[columnIndex+1]); i++ {
		newNeuron[i] = DefaultAxonWeight
	}

	// Give the previous column a weight for this new neuron
	for i := 0; i < len(nn[columnIndex-1]); i++ {
		nn[columnIndex-1][i] = append(nn[columnIndex-1][i], DefaultAxonWeight)
	}

	nn[columnIndex] = append(nn[columnIndex], newNeuron)
}

/**
 * Remove the column between our output column and the column two indexes prior.
 */
func (nn_p *RectifierNetwork) removeColumn(cOpt_p *RectifierColumnGenerationOptions) *RectifierNetwork {

	nn := *nn_p
	cOpt := *cOpt_p

	i := len(nn) - 2

	oldLength := len(nn[i])

	nn = append(nn[:i], nn[i+1:]...)

	i--

	// Scrap the weights which we don't have anymore,
	// if the out column is smaller than the column
	// which used to precede it.
	if oldLength > len(nn[i+1]) {
		for j := 0; j < len(nn[i]); j++ {
			nn[i][j] = nn[i][j][:len(nn[i+1])]
		}
		// Add more default weights, if the out column
		// is larger.
	} else if oldLength < len(nn[i+1]) {
		for j := 0; j < len(nn[i]); j++ {
			for k := len(nn[i]); k < len(nn[i+1]); k++ {
				nn[i][j] = append(nn[i][j], cOpt.DefaultAxonWeight)
			}
		}
	}

	return &nn
}

/**
 * Add a column between our output column and the column immediately prior.
 */
func (nn_p *RectifierNetwork) addColumn(cOpt_p *RectifierColumnGenerationOptions) *RectifierNetwork {

	nn := *nn_p
	cOpt := *cOpt_p

	i := len(nn) - 1

	// The outColumn is just an array of no weights,
	// so we don't need to worry about copying the
	// old column over.
	outColumn := make([]RectifierNeuron, len(nn[i]))

	// Add a bunch of default weights to
	// the current out column, converting it into
	// a regular column.
	for j := 0; j < len(nn[i]); j++ {
		nn[i][j] = make(RectifierNeuron, len(nn[i]))
		for k := 0; k < len(nn[i]); k++ {
			nn[i][j][k] = cOpt.DefaultAxonWeight
		}
	}

	// Add the out column back onto the network
	nn = append(nn, outColumn)

	newSize := (rand.Intn(cOpt.MaxSize-cOpt.MinSize) + cOpt.MinSize) - len(nn[i])

	// Add some more neurons to the old outColumn.
	for j := 0; j < newSize; j++ {
		nn.addNeuron(i, cOpt.DefaultAxonWeight)
	}

	return &nn
}

/**
 * Swap two Axons which start from these neurons indexes.
 */
func (nn_p *RectifierNetwork) swapWeights(columnIndex, neuronIndex, axonIndex1, axonIndex2 int) {

	neuron := (*nn_p)[columnIndex][neuronIndex]

	weight1 := neuron[axonIndex1]
	neuron[axonIndex1] = neuron[axonIndex2]
	neuron[axonIndex2] = weight1
}