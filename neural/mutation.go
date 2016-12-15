package neural

import (
	"math/rand"
)

func (genOpt NetworkGenerationOptions) Mutate(n *Network) *Network {
	return n.MutateOpts(genOpt.NetworkMutationOptions)
}

/**
 * Modify a float to some float close to it in value.
 * TODO: Give this a better distribution of randomness.
 */
func mutateFloat(toMutate float64, opt FloatMutationOptions) float64 {
	if rand.Float64() >= opt.MutChance {
		return toMutate
	}
	if rand.Float64() <= opt.ZeroOutChance {
		return 0.0
	}

	out := (opt.MutMagnitude * float64(rand.Intn((opt.MutRange*2)+1))) + toMutate
	return out - (opt.MutMagnitude * float64(opt.MutRange))
}

/**
 * Mutate this network
 */
func (nn *Network) MutateOpts(mOpt NetworkMutationOptions) *Network {

	nn.Body.Mutate(mOpt)

	if rand.Float64() < mOpt.ActivatorMutationChance {
		nn.Activator = MutateActivator(mOpt.ActivatorOptions)
	}

	return nn
}

func (nn *Network) Mutate() {
	nn.MutateOpts(ngo.NetworkMutationOptions)
}

// All activators are currently weighted the same
// in this mutation. A future implementation could
// also take in how highly each activator should be
// weighed, and could avoid mutating into the current
// activator (although this only has the effect of
// slightly reducing real mutation chance)
/**
 * Mutate an activator function.
 */
func MutateActivator(mOpt ActivatorMutationOptions) ActivatorFunc {
	return mOpt[rand.Intn(len(mOpt))]
}

/**
 * Mutate this network body.
 */
func (b *Body) Mutate(mOpt NetworkMutationOptions) {

	bv := *b

	if rand.Float64() < mOpt.ColumnRemovalChance {
		// We currently only remove the len-1th column
		// We can't remove a column if we have just one
		// hidden layer or the output column to remove
		if len(bv) > 2 {
			b.removeColumn(mOpt.ColumnOptions)
		}
	}

	if rand.Float64() < mOpt.ColumnAdditionChance {
		// We currently only add a column in the len-1th space
		// In the future a check against some maximum column count
		// should be here
		b.addColumn(mOpt.ColumnOptions)
	}

	for i := 0; i < len(bv)-1; i++ {

		col := bv[i]

		// Swap two axons connecting this column
		// to the next column
		if rand.Float64() < mOpt.WeightSwapChance {

			// We can't make a meaningful swap
			// if the column only has one neuron
			if len(col) > 1 {
				// Get our neuron index
				neuronIndex := rand.Intn(len(col))

				// Get our first axon index
				axonIndex1 := rand.Intn(len(bv[i+1]))

				// Get our second axon index
				axonIndex2 := rand.Intn(len(bv[i+1]))
				// If our second random neuron is same
				// as our first, we interpret that as
				// taking the neuron just following our first
				// as our second
				if axonIndex2 == axonIndex1 {
					axonIndex2 = (axonIndex2 + 1) % len(bv[i+1])
				}

				b.swapWeights(i, neuronIndex, axonIndex1, axonIndex2)
			}
		}

		if i != 0 {
			if rand.Float64() < mOpt.NeuronAdditionChance {
				// In the future a check on some max length
				// for a column should exist.
				b.addNeuron(i, (mOpt.ColumnOptions).DefaultAxonWeight)
			}

			if rand.Float64() < mOpt.NeuronReplacementChance {
				neuronIndex := rand.Intn(len(col))
				b.replaceNeuron(i, neuronIndex, (mOpt.ColumnOptions).DefaultAxonWeight)
			}
		}
	}

	// Mutate individual neurons
	for x := 0; x < len(bv); x++ {
		for y := 0; y < len(bv[x]); y++ {
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
				b.mutateNeuron(x, y, mOpt.WeightOptions)
			}
		}
	}
}

func (b *Body) mutateNeuron(columnIndex, neuronIndex int, wOpt FloatMutationOptions) {

	bv := *b

	newNeuron := make(Neuron, len(bv[columnIndex][neuronIndex]))
	for i, weight := range bv[columnIndex][neuronIndex] {
		newNeuron[i] = mutateFloat(weight, wOpt)
	}

	bv[columnIndex][neuronIndex] = newNeuron
	*b = bv
}

/**
 * Removes a neuron from an index and places a new neuron there in its place.
 * Effectively resetNeuron + addNeuron, if addNeuron took an index.
 */
func (b *Body) replaceNeuron(columnIndex, neuronIndex int, DefaultAxonWeight float64) {

	bv := *b

	for i := 0; i < len(bv[columnIndex+1]); i++ {
		bv[columnIndex][neuronIndex][i] = DefaultAxonWeight
	}

	*b = bv
}

/**
 * Add a neuron to the end of a column.
 */
func (b *Body) addNeuron(columnIndex int, DefaultAxonWeight float64) {

	bv := *b

	newNeuron := make(Neuron, len(bv[columnIndex+1]))

	// Set this new neuron's weights for the next column
	for i := 0; i < len(bv[columnIndex+1]); i++ {
		newNeuron[i] = DefaultAxonWeight
	}

	// Give the previous column a weight for this new neuron
	for i := 0; i < len(bv[columnIndex-1]); i++ {
		bv[columnIndex-1][i] = append(bv[columnIndex-1][i], DefaultAxonWeight)
	}

	bv[columnIndex] = append(bv[columnIndex], newNeuron)

	*b = bv

}

/**
 * Remove the column between our output column and the column two indexes prior.
 */
func (b *Body) removeColumn(cOpt ColumnGenerationOptions) {

	bv := *b

	i := len(bv) - 2

	oldLength := len(bv[i])

	bv = append(bv[:i], bv[i+1:]...)

	i--

	// Scrap the weights which we don't have anymore,
	// if the out column is smaller than the column
	// which used to precede it.
	if oldLength > len(bv[i+1]) {
		for j := 0; j < len(bv[i]); j++ {
			bv[i][j] = bv[i][j][:len(bv[i+1])]
		}
		// Add more default weights, if the out column
		// is larger.
	} else if oldLength < len(bv[i+1]) {
		for j := 0; j < len(bv[i]); j++ {
			for k := len(bv[i]); k < len(bv[i+1]); k++ {
				bv[i][j] = append(bv[i][j], cOpt.DefaultAxonWeight)
			}
		}
	}

	*b = bv
}

/**
 * Add a column between our output column and the column immediately prior.
 */
func (b *Body) addColumn(cOpt ColumnGenerationOptions) {

	bv := *b

	i := len(bv) - 1

	// The outColumn is just an array of no weights,
	// so we don't need to worry about copying the
	// old column over.
	outColumn := make([]Neuron, len(bv[i]))

	// Add a bunch of default weights to
	// the current out column, converting it into
	// a regular column.
	for j := 0; j < len(bv[i]); j++ {
		bv[i][j] = make(Neuron, len(bv[i]))
		for k := 0; k < len(bv[i]); k++ {
			bv[i][j][k] = cOpt.DefaultAxonWeight
		}
	}

	// Add the out column back onto the network
	bv = append(bv, outColumn)

	newSize := (rand.Intn(cOpt.MaxSize-cOpt.MinSize) + cOpt.MinSize) - len(bv[i])

	*b = bv

	// Add some more neurons to the old outColumn.
	for j := 0; j < newSize; j++ {
		b.addNeuron(i, cOpt.DefaultAxonWeight)
	}
}

/**
 * Swap two Axons which start from these neurons indexes.
 */
func (nn_p *Body) swapWeights(columnIndex, neuronIndex, axonIndex1, axonIndex2 int) {

	neuron := (*nn_p)[columnIndex][neuronIndex]

	weight1 := neuron[axonIndex1]
	neuron[axonIndex1] = neuron[axonIndex2]
	neuron[axonIndex2] = weight1
}
