package neural

import (
	"math/rand"
	"fmt"
)
type NetworkMutationOptions struct {
	neuronOptions *NeuronMutationOptions
	columnOptions *ColumnGenerationOptions
	// per column
	neuronReplacementChance float64 
	neuronAdditionChance float64 
	axonRemovalChance float64 
	axonAdditionChance float64 
	axonSwapChance float64
	// per network
	columnRemovalChance float64 
	columnAdditionChance float64
	neuronMutationChance float64 
}
type NeuronMutationOptions struct {
	thresholdOptions FloatMutationOptions
	weightOptions FloatMutationOptions
}
type FloatMutationOptions struct {
	mutChance float64
	mutMagnitude float64
	mutRange int
}
type NetworkGenerationOptions struct {
	NetworkMutationOptions
	minColumns int
	maxColumns int
	inputs int
	outputs int 
	baseMutations int
}
type ColumnGenerationOptions struct {
	minSize int
	maxSize int
	neuronOptions *NeuronGenerationOptions 
}
type NeuronGenerationOptions struct {
	minAxons int
	maxAxons int
	defaultThreshold float64
	defaultAxonWeight float64
}

/**
 * Modify a float to some float close to it in value. 
 */
func mutateFloat(toMutate float64, opt FloatMutationOptions) float64 {
	if rand.Float64() >= opt.mutChance {
		return toMutate
	}
	
	out := (opt.mutMagnitude * float64(rand.Intn(opt.mutRange*2))) + toMutate
	return out - (opt.mutMagnitude * float64(opt.mutRange))
}

/**
 * Mutate this neuron.
 */ 
func (n *Neuron) mutate(mOpt_p *NeuronMutationOptions) Neuron {
	mOpt := *mOpt_p

	newThreshold := mutateFloat(n.threshold, mOpt.thresholdOptions)

	newWeights := map[int]float64{}
	for i,weight := range n.weights {
		newWeights[i] = mutateFloat(weight, mOpt.weightOptions)
	}

	return Neuron{
			outputs:n.outputs, 
			threshold:newThreshold,
			weights:newWeights,
			val:0,
		}
} 

/**
 * Mutate this network.
 */
func (nn *Network) Mutate(mOpt_p *NetworkMutationOptions) *Network {
	
	mOpt := *mOpt_p

	newNetwork := nn.copy()

	// We store this and continually shift it
	// a digit to the left in case rand.Float64()
	// is more costly than a multiplication,
	// a subtraction, and a floor call. 
	// This might be wrong. 
	randCh := make(chan float64)
	go RefreshingRand(randCh)

	if <-randCh < mOpt.columnRemovalChance {
		// We currently only remove the len-1th column
		// We can't remove a column if we have just one
		// hidden layer or the output column to remove
		if len(newNetwork) > 2 {
			newNetwork.removeColumn(mOpt.columnOptions)
		}
	}

	if <-randCh < mOpt.columnAdditionChance {
		// We currently only add a column in the len-1th space
		// In the future a check against some maximum column count
		// should be here
		newNetwork.addColumn(mOpt.columnOptions)
	}

	for i := 0; i < len(newNetwork)-1; i++ {

		col := newNetwork[i]

		// Swap two axons connecting this column
		// to the next column
		if <-randCh < mOpt.axonSwapChance {

			// We can't make a meaningful swap
			// if the column only has one neuron
			if len(col) > 1 {
				// Get our first neuron, axon pair
				neuronIndex1 := rand.Intn(len(col))

				// Get our second neuron, axon pair
				neuronIndex2 := rand.Intn(len(col))
				// If our second random neuron is same
				// as our first, we interpret that as
				// taking the neuron just following our first
				// as our second
				if neuronIndex2 == neuronIndex1 {
					neuronIndex2 = (neuronIndex2 + 1) % len(col)
				}

				// We can't swap axons on neurons who
				// don't have output axons.
				if len(col[neuronIndex1].outputs) > 0 && 
				   len(col[neuronIndex2].outputs) > 0 {
					newNetwork.swapAxons(i, neuronIndex1, neuronIndex2)
				}
			}
		}

		// Remove an axon connecting this column
		// to the next column
		if <-randCh < mOpt.axonRemovalChance {
			// We don't want to remove the only
			// axon a neuron has. 
			neuronIndex := rand.Intn(len(col))
			if len(col[neuronIndex].outputs) > 1 {
				newNetwork.removeAxon(i, neuronIndex)
			}	 
		}

		// Add an additional axon from this column
		// to the next column
		if <-randCh < mOpt.axonAdditionChance {
			// We can't add another axon if we picked
			// a neuron which is already connected to
			// every neuron in the next column.
			//
			// We could retry to find another neuron.
			// We aren't doing that.
			// It would lead to looping over all neurons
			// until we found one that was valid
			// or ejecting once we couldn't
			// and skipping an iteration of mutation
			// is much better for us performance-wise.

			axonWeight := (*((mOpt.columnOptions).neuronOptions)).defaultAxonWeight

			neuronIndex := rand.Intn(len(col))
			if len(newNetwork[i+1]) != len(col[neuronIndex].outputs) {
				newNetwork.addAxon(i, neuronIndex, axonWeight)
			}
		}

		if i != 0 {
			if <-randCh < mOpt.neuronAdditionChance {
				// In the future a check on some max length
				// for a column should exist.
				newNetwork.addNeuron(i, (*mOpt.columnOptions).neuronOptions)
			}

			if <-randCh < mOpt.neuronReplacementChance {
				// We can't remove a column's only neuron
				if len(col) != 1 {
					neuronIndex := rand.Intn(len(col))
					newNetwork.replaceNeuron(i, neuronIndex, (*mOpt.columnOptions).neuronOptions)
				}
			}
		}
	}

	// Mutate individual neurons
	for x,column := range newNetwork {
		for y, neuron := range column {

			// For performance in the future,
			// We should replace having x * y random calls
			// with picking an average number of neurons
			// to always mutate where the average is 
			// calculated by x * y * mutChance. Then
			// we'd make x * y * mutChance random index
			// to-mutate determinations. 
			//
			// This also applies to rand.Float64() calls
			// above. 
			if <-randCh < mOpt.neuronMutationChance {
				newNetwork[x][y] = neuron.mutate(mOpt.neuronOptions)
			}
		}
	}

	return &newNetwork
}

/**
 * Create a copy of this network's output column.
 * Used when adding and removing columns next to the output column.
 */ 
func (nn_p *Network) copyOutColumn() []Neuron {

	nn := *nn_p

	i := len(nn) - 1 

	// Create a new outColumn
	outColumn := []Neuron{}
	for j := 0; j < len(nn[i]); j++ {
		outColumn = append(outColumn, 
						    Neuron{
						        threshold: nn[i][j].threshold,
						        weights: make(map[int]float64),
								outputs: make(map[int]bool),
								val: 0,
						    })
	}
	return outColumn
}

/**
 * Remove a random neuron from a column.
 */
func (nn_p *Network) removeNeuron(columnIndex int, neuronIndex int) {

	nn := *nn_p

	neuron := nn[columnIndex][neuronIndex]
	prevCol := nn[columnIndex-1]

	// Axons connecting to this neuron
	for index := range neuron.weights {
		delete(prevCol[index].outputs, neuronIndex)
	}

	// Axons leaving from this neuron
	if columnIndex != (len(nn)-1) {
		nextCol := nn[columnIndex+1]
		for index := range neuron.outputs {
			delete(nextCol[index].weights, neuronIndex)
		}
	}

	col := nn[columnIndex]

	// Remove this neuron from our column
	nn[columnIndex] = append(col[:neuronIndex],col[neuronIndex+1:]...)
}

/**
 * Removes a neuron from an index and places a new neuron there in its place.
 * Effectively removeNeuron + addNeuron, if removeNeuron didn't shrink
 * and addNeuron took an index.
 */
func (nn_p *Network) replaceNeuron(columnIndex int, neuronIndex int, nOpt_p *NeuronGenerationOptions) {

	nOpt := *nOpt_p
	nn := *nn_p

	neuron := nn[columnIndex][neuronIndex]
	prevCol := nn[columnIndex-1]
	nextCol := nn[columnIndex+1]

	// Delete Axons connecting to this neuron
	for index := range neuron.weights {
		delete(prevCol[index].outputs, neuronIndex)
	}

	// Delete Axons leaving from this neuron
	for index := range neuron.outputs {
		delete(nextCol[index].weights, neuronIndex)
	}

	nn[columnIndex][neuronIndex].threshold = nOpt.defaultThreshold

	// Create new Axons connecting to this neuron
	axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxonBack(columnIndex, neuronIndex, nOpt.defaultAxonWeight)
	}
	// Create new Axons leaving from this neuron
	axonCount = rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxon(columnIndex, neuronIndex, nOpt.defaultAxonWeight)
	}
}

/**
 * Add a neuron to the end of a column.
 */
func (nn_p *Network) addNeuron(columnIndex int, nOpt_p *NeuronGenerationOptions) {

	nn := *nn_p
	nOpt := *nOpt_p

	nn[columnIndex] = append(nn[columnIndex], 
		Neuron{
			threshold:nOpt.defaultThreshold,
			outputs:make(map[int]bool),
			weights:make(map[int]float64),
			val:0,
		})

	// Axons connecting to this neuron
	axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxonBack(columnIndex, len(nn[columnIndex])-1, nOpt.defaultAxonWeight)
	}
	// Axons leaving from this neuron
	axonCount = rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxon(columnIndex, len(nn[columnIndex])-1, nOpt.defaultAxonWeight)
	}
}

/**
 * Remove the column between our output column and the column two indexes prior.
 */
func (nn_p *Network) removeColumn(cOpt_p *ColumnGenerationOptions) {

	nn := *nn_p
	cOpt := *cOpt_p
	nOpt := *(cOpt.neuronOptions)

	i := len(nn) - 2

	outColumn := nn.copyOutColumn()

	// Remove neuron will deal with
	// the axons that need to be disconnected
	// from this column
	for len(nn[i]) > 0 {
		nn.removeNeuron(i, 0)
	}

	nn = append(nn[:i],outColumn)

	i--

	// Add in some random connections from the 
	// new final column to the output column.
	for j := 0; j < len(nn[i]); j++ {
		axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
		for k := 0; k < axonCount; k++ {
			nn.addAxon(i,j, nOpt.defaultAxonWeight)
		}
	}
}

/**
 * Add a column between our output column and the column immediately prior.
 */
func (nn_p *Network) addColumn(cOpt_p *ColumnGenerationOptions) *Network {

	nn := *nn_p
	cOpt := *cOpt_p

	i := len(nn) - 1 

	outColumn := nn.copyOutColumn()

	// Replace the current outColumn with
	// an empty column, disconnecting all
	// neurons from the previous column
	// in the process.
	for len(nn[i]) > 0 {
		nn.removeNeuron(i, 0)
	}

	// Place the new outColumn after the
	// empty column.
	nn = append(nn, outColumn)

	newSize := rand.Intn(cOpt.maxSize-cOpt.minSize) + cOpt.minSize

	// Repopulate the empty column with new neurons
	// and connections to the bordering columns.
	for j := 0; j < newSize; j++ {
		nn.addNeuron(i, cOpt.neuronOptions)
	}

	return &nn
}

/**
 * Swap two Axons which start from these neurons indexes.
 */
func (nn_p *Network) swapAxons(columnIndex int, neuronIndex1 int, neuronIndex2 int) {

	nn := *nn_p

	neuron1 := nn[columnIndex][neuronIndex1]
	neuron2 := nn[columnIndex][neuronIndex2]

	// This list generation could be removed
	// with a data structure that kept track
	// of a key list along with a map, or
	// otherwise provided random-element
	// access to a map. 
	axonList1 := KeySet(neuron1.outputs)
	axonIndex1 := axonList1[rand.Intn(len(axonList1))]

	nn.print()
	fmt.Println(axonIndex1, columnIndex, neuronIndex1)

	axon1 := nn[columnIndex+1][axonIndex1]
	
	axonList2 := KeySet(neuron2.outputs)
	axonIndex2 := axonList2[rand.Intn(len(axonList2))]
	axon2 := nn[columnIndex+1][axonIndex2]

	// Swap endpoints
	// We also swap weights here.
	// An alternative swap would
	// maintain weights within a neuron.
	weight1 := axon1.weights[neuronIndex1]
	axon1.weights[neuronIndex2] = axon2.weights[neuronIndex2]
	axon2.weights[neuronIndex1] = weight1
	// Clean up old endpoints
	delete(axon1.weights, neuronIndex1)
	delete(axon2.weights, neuronIndex2)

	// Swap sendpoints
	neuron1.outputs[axonIndex2] = true
	neuron2.outputs[axonIndex1] = true
	// Clean up old sendpoints
	delete(neuron1.outputs, axonIndex1)
	delete(neuron2.outputs, axonIndex2)
}

/**
 * Remove a random Axon which starts from this neuron's index.
 */
func (nn *Network) removeAxon(columnIndex int, neuronIndex int) {

	neuron := (*nn)[columnIndex][neuronIndex]

	axonIndex := KeySet(neuron.outputs)[rand.Intn(len(neuron.outputs))]
	axon := (*nn)[columnIndex+1][axonIndex]

	// Delete the endpoint
	delete(axon.weights, neuronIndex)

	// Delete the sendpoint 
	delete(neuron.outputs, axonIndex)
}

/**
 * Add a random Axon which starts from this neuron's index.
 */
func (nn_p *Network) addAxon(columnIndex int, neuronIndex int, defaultAxonWeight float64) {

	nn := *nn_p

	nextCol := nn[columnIndex+1]
	neuron := nn[columnIndex][neuronIndex]

	if len(neuron.outputs) >= len(nextCol) {
		return
	}

	// Choices is a map of valid choices
	// for our new axon to connect to
	choices := make(map[int]bool)

	// We start off with all indexes
	// in the next column as valid choices
	for i := range nextCol {
		choices[i] = true
	}

	// We then remove all indexes which
	// already exist in our neuron's outputs
	for choice := range neuron.outputs {
		delete(choices, choice)
	}

	choice := KeySet(choices)[rand.Intn(len(choices))]

	nextCol[choice].weights[neuronIndex] = defaultAxonWeight

	neuron.outputs[choice] = true
}

/**
 * Add a random Axon which ends at this neuron's index.
 */
func (nn_p *Network) addAxonBack(columnIndex int, neuronIndex int, defaultAxonWeight float64) {
	
	nn := *nn_p

	nextCol := nn[columnIndex-1]
	neuron := nn[columnIndex][neuronIndex]

	if len(neuron.weights) == len(nextCol) {
		return
	}

	// Choices is a map of valid choices
	// for our new axon to connect to
	choices := make(map[int]bool)

	// We start off with all indexes
	// in the next column as valid choices
	for i := range nextCol {
		choices[i] = true
	}

	// We then remove all indexes which
	// already exist in our neuron's weights
	for choice := range neuron.weights {
		delete(choices, choice)
	}

	choice := KeySet(choices)[rand.Intn(len(choices))]

	nextCol[choice].outputs[neuronIndex] = true
	neuron.weights[choice] = defaultAxonWeight
}