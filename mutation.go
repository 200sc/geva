package neural

import (
	"rand"
)

type NetworkMutationOptions struct {
	neuronOptions *NeuronMutationOptions
	columnOptions *ColumnGenerationOptions
	// per column
	neuronRemovalChance float64 
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
func mutateFloat(toMutate float64, opt FloatMutationOptions) float {
	if rand.Float64() >= opt.mutChance {
		return toMutate
	}
	
	out := (opt.mutMagnitude * rand.Intn(opt.mutRange*2)) + toMutate
	return out - (opt.mutMagnitude * opt.mutRange)
}

/**
 * Mutate this neuron.
 */ 
func (n *Neuron) mutate(mOpt_p *NeuronMutationOptions) Neuron {
	mOpt := *mOpt_p

	newThreshold := mutateFloat(n.threshold, mOpt.thresholdOptions)

	newWeights := []float
	for _,weight := range n.weights {
		newWeights = append(weights, mutateFloat(weight, mOpt.weightOptions))
	}

	return Neuron{n.outputs, newThreshold, newWeights}
} 

/**
 * Mutate this network.
 * Note how this is the only public function (as of this comment).
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
		if len(newNetwork.columns) > 2 {
			newNetwork.removeColumn(mOpt.columnOptions)
		}
	}

	if <-randCh < mOpt.columnAdditionChance {
		// We currently only add a column in the len-1th space
		// In the future a check against some maximum column count
		// should be here
		newNetwork.addColumn(mOpt.columnOptions)
	}

	for i := range len(newNetwork.columns-1) {

		col := newNetwork.columns[i]

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
				// taking the neuron just prior in order
				// to our first neuron as our second
				if neuronIndex2 == neuronIndex1 {
					neuronIndex2 = (neuronIndex2 - 1) % len(col)
				}

				newNetwork.swapAxons(i, neuronIndex1, neuronIndex2)
			}
		}

		// Remove an axon connecting this column
		// to the next column
		if <-randCh < mOpt.axonRemovalChance {
			// We don't want to remove the only
			// axon a neuron has. 
			neuronIndex := rand.Intn(len(col))
			if len(col[neuronIndex].outputs) != 1 {
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

			axonWeight = *(*(mOpt.columnOptions).neuronOptions).defaultAxonWeight

			neuronIndex := rand.Intn(len(col))
			if len(newNetwork.columns[i+1]) != len(col[neuronIndex].outputs) {
				newNetwork.addAxon(i, neuronIndex, axonWeight)
			}
		}

		if i != 0 {
			if <-randCh < mOpt.neuronAdditionChance {
				// In the future a check on some max length
				// for a column should exist.
				newNetwork.addNeuron(i, (*mOpt.columnOptions).neuronOptions)
			}

			if <-randCh < mOpt.neuronRemovalChance {
				// We can't remove a column's only neuron
				if len(col) != 1 {
					neuronIndex := rand.Intn(len(col))
					newNetwork.removeNeuron(i, neuronIndex)
				}
			}
		}
	}

	// Mutate individual neurons
	for x,column := range nn.columns {
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
				newNetwork.columns[x][y] = neuron.mutate(mOpt.neuronOptions)
			}
		}
	}

	return newNetwork
}

/**
 * Create a copy of this network's output column.
 * Used when adding and removing columns next to the output column.
 */ 
func (nn *Network) copyOutColumn() (outColumn []Neuron) {

	i := len(nn.columns) - 1 

	// Create a new outColumn
	outColumn := []Neuron{}
	for j := 0; j < len(nn.columns[i]); j++ {
		outColumn = append(outColumn, 
						    Neuron{
						    	// We want no outputs on our new
						    	// output column, and our weights
						    	// will be determined by addAxon calls.  
						    	// That just leaves threshold.
						        threshold: nn.columns[i][j].threshold
						    })
	}
}

/**
 * Remove a random neuron from a column.
 */
func (nn *Network) removeNeuron(columnIndex int, neuronIndex int) {

	neuron := nn.columns[columnIndex][neuronIndex]
	prevCol := nn.columns[columnIndex-1]
	nextCol := nn.columns[columnIndex+1]

	// Axons connecting to this neuron
	for index := range neuron.weights {
		delete(prevCol[index].outputs, neuronIndex)
	}

	// Axons leaving from this neuron
	for index := range neuron.outputs {
		delete(nextCol[index].weights, neuronIndex)
	}

	col := nn.columns[columnIndex]

	// Remove this neuron from our column
	nn.columns[columnIndex] = col[:neuronIndex] + col[neuronIndex+1:]
}

/**
 * Add a neuron to the end of a column.
 */
func (nn *Network) addNeuron(columnIndex int, nOpt_p *NeuronGenerationOptions) {

	nOpt := *nOpt_p

	nn.columns[columnIndex] = append(nn.columns[columnIndex], Neuron{threshold:nOpt.defaultThreshold})

	// Axons connecting to this neuron
	axonCount := rand.Intn(cOpt.maxAxns-cOpt.minAxons) + cOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxonBack(columnIndex, len(nn.columns[columnIndex])-1, nOpt.defaultAxonWeight)
	}
	// Axons leaving from this neuron
	axonCount = rand.Intn(cOpt.maxAxns-cOpt.minAxons) + cOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxon(columnIndex, len(nn.columns[columnIndex])-1, nOpt.defaultAxonWeight)
	}
}

/**
 * Remove the column between our output column and the column two indexes prior.
 */
func (nn *Network) removeColumn(cOpt_p *ColumnGenerationOptions) {

	cOpt := *cOpt_p
	nOpt := *(cOpt.neuronOptions)

	i := len(nn.columns) - 2

	outColumn := nn.copyOutColumn()

	// Remove column will deal with
	// the axons that need to be disconnected
	// from this column
	for len(nn.columns[i]) > 0 {
		nn.removeNeuron(i)
	}

	nn.columns = nn.columns[:i] + outColumn

	i--

	// Add in some random connections from the 
	// new final column to the output column.
	for j := 0; j < len(nn.columns[i]); j++ {
		axonCount := rand.Intn(cOpt.maxAxns-cOpt.minAxons) + cOpt.minAxons
		for k := 0; k < axonCount; k++ {
			nn.addAxon(i,j, nOpt.defaultAxonWeight)
		}
	}
}

/**
 * Add a column between our output column and the column immediately prior.
 */
func (nn *Network) addColumn(cOpt_p *ColumnGenerationOptions) {

	cOpt := *cOpt_p

	i := len(nn.columns) - 1 

	outColumn := nn.copyOutColumn()

	// Replace the current outColumn with
	// an empty column, disconnecting all
	// neurons from the previous column
	// in the process.
	for len(nn.columns[i]) > 0 {
		nn.removeNeuron(i)
	}

	// Place the new outColumn after the
	// empty column.
	nn.columns = append(nn.columns, outColumn)

	newSize := rand.Intn(cOpt.maxSize-cOpt.minSize) + cOpt.minSize

	// Repopulate the empty column with new neurons
	// and connections to the bordering columns.
	for j := 0; j < newSize; j++ {
		nn.addNeuron(i)
	}
}

/**
 * Swap two Axons which start from this neuron's index.
 */
func (nn *Network) swapAxons(columnIndex int, neuronIndex1 int, neuronIndex2 int) {

	neuron1 := nn.columns[columnIndex][neuronIndex1]
	neuron2 := nn.columns[columnIndex][neuronIndex2]

	// This list generation could be removed
	// with a data structure that kept track
	// of a key list along with a map, or
	// otherwise provided random-element
	// access to a map. 
	axonList1 := KeySet(neuron1.outputs)
	axonIndex1 := axonList1[rand.Intn(len(axonList1))]
	axon1 := newNetwork.columns[i+1][axonIndex1]
	
	axonList2 := KeySet(neuron2.outputs)
	axonIndex2 := axonList2[rand.Intn(len(axonList2))]
	axon2 := newNetwork.columns[i+1][axonIndex2]

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

	neuron := nn[columnIndex][neuronIndex]

	axonIndex := neuron.outputs[rand.Intn(len(neuron.ouputs))]
	axon := nn.columns[columnIndex+1][axonIndex]

	// Delete the endpoint
	delete(axon.weights, neuronIndex)

	// Delete the sendpoint 
	delete(neuron.outputs, axonIndex)
}

/**
 * Add a random Axon which starts from this neuron's index.
 */
func (nn *Network) addAxon(columnIndex int, neuronIndex int, defaultAxonWeight int) {

	nextCol := nn.columns[columnIndex+1]
	neuron := nn.columns[columnIndex][neuronIndex]

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
	// MAGIC NUMBER ALERT
	// DEFAULT WEIGHT IS 1 BUT THAT ISN'T DETAILED
	nextCol[choice].weights[neuronIndex] = defaultAxonWeight
	neuron.outputs[choice] = true
}

/**
 * Add a random Axon which ends at this neuron's index.
 */
func (nn *Network) addAxonBack(columnIndex int, neuronIndex int, defaultAxonWeight int) {
	
	nextCol := nn.columns[columnIndex-1]
	neuron := nn.columns[columnIndex][neuronIndex]

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
	// MAGIC NUMBER ALERT
	// DEFAULT WEIGHT IS 1 BUT THAT ISN'T DETAILED
	nextCol[choice].outputs[neuronIndex] = true
	neuron.weights[choice] = defaultAxonWeight
}