package neural

import (
	"math/rand"
)

type PerceptronNetworkMutationOptions struct {
	neuronOptions *PerceptronMutationOptions
	ColumnOptions *PerceptronColumnGenerationOptions
	// per column
	NeuronReplacementChance float64
	NeuronAdditionChance    float64
	axonRemovalChance       float64
	axonAdditionChance      float64
	axonSwapChance          float64
	// per network
	ColumnRemovalChance  float64
	ColumnAdditionChance float64
	NeuronMutationChance float64
}
type PerceptronNetworkGenerationOptions struct {
	PerceptronNetworkMutationOptions
	MinColumns    int
	MaxColumns    int
	Inputs        int
	Outputs       int
	BaseMutations int
}
type PerceptronMutationOptions struct {
	thresholdOptions FloatMutationOptions
	WeightOptions    FloatMutationOptions
}
type PerceptronGenerationOptions struct {
	minAxons          int
	maxAxons          int
	defaultThreshold  float64
	DefaultAxonWeight float64
}
type PerceptronColumnGenerationOptions struct {
	MinSize       int
	MaxSize       int
	neuronOptions *PerceptronGenerationOptions
}

func (genOpt PerceptronNetworkGenerationOptions) Generate() Network {
	return *GeneratePerceptronNetwork(&genOpt)
}

func (genOpt PerceptronNetworkGenerationOptions) Mutate(n Network) Network {
	return n.(PerceptronNetwork).Mutate(&(genOpt.PerceptronNetworkMutationOptions))
}

/**
 * Modify a float to some float close to it in value.
 * TODO: Give this a better distribution of randomness.
 */
func mutateFloat(toMutate float64, opt FloatMutationOptions) float64 {
	if rand.Float64() >= opt.MutChance {
		return toMutate
	}

	out := (opt.MutMagnitude * float64(rand.Intn(opt.MutRange*2))) + toMutate
	return out - (opt.MutMagnitude * float64(opt.MutRange))
}

/**
 * Mutate this Perceptron.
 */
func (n *Perceptron) mutate(mOpt_p *PerceptronMutationOptions) Perceptron {
	mOpt := *mOpt_p

	newThreshold := mutateFloat(n.threshold, mOpt.thresholdOptions)

	newWeights := map[int]float64{}
	for i, weight := range n.weights {
		newWeights[i] = mutateFloat(weight, mOpt.WeightOptions)
	}

	return Perceptron{
		Outputs:   n.Outputs,
		threshold: newThreshold,
		weights:   newWeights,
	}
}

/**
 * Mutate this network.
 */
func (nn PerceptronNetwork) Mutate(mOpt_p *PerceptronNetworkMutationOptions) *PerceptronNetwork {

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
		if rand.Float64() < mOpt.axonSwapChance {

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
				if len(col[neuronIndex1].Outputs) > 0 &&
					len(col[neuronIndex2].Outputs) > 0 {
					newNetwork.swapAxons(i, neuronIndex1, neuronIndex2)
				}
			}
		}

		// Remove an axon connecting this column
		// to the next column
		if rand.Float64() < mOpt.axonRemovalChance {
			// We don't want to remove the only
			// axon a neuron has.
			neuronIndex := rand.Intn(len(col))
			if len(col[neuronIndex].Outputs) > 1 {
				newNetwork.removeAxon(i, neuronIndex)
			}
		}

		// Add an additional axon from this column
		// to the next column
		if rand.Float64() < mOpt.axonAdditionChance {
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

			axonWeight := (*((mOpt.ColumnOptions).neuronOptions)).DefaultAxonWeight

			neuronIndex := rand.Intn(len(col))
			if len(newNetwork[i+1]) != len(col[neuronIndex].Outputs) {
				newNetwork.addAxon(i, neuronIndex, axonWeight)
			}
		}

		if i != 0 {
			if rand.Float64() < mOpt.NeuronAdditionChance {
				// In the future a check on some max length
				// for a column should exist.
				newNetwork.addNeuron(i, (*mOpt.ColumnOptions).neuronOptions)
			}

			if rand.Float64() < mOpt.NeuronReplacementChance {
				// We can't remove a column's only neuron
				if len(col) > 1 {
					neuronIndex := rand.Intn(len(col))
					newNetwork.replaceNeuron(i, neuronIndex, (*mOpt.ColumnOptions).neuronOptions)
				}
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
func (nn_p *PerceptronNetwork) copyOutColumn() []Perceptron {

	nn := *nn_p

	i := len(nn) - 1

	// Create a new outColumn
	outColumn := []Perceptron{}
	for j := 0; j < len(nn[i]); j++ {
		outColumn = append(outColumn,
			Perceptron{
				threshold: nn[i][j].threshold,
				weights:   make(map[int]float64),
				Outputs:   make(map[int]bool),
			})
	}
	return outColumn
}

/**
 * Remove all connections from a neuron
 */
func (nn_p *PerceptronNetwork) resetNeuron(columnIndex int, neuronIndex int) {

	nn := *nn_p

	neuron := nn[columnIndex][neuronIndex]
	prevCol := nn[columnIndex-1]

	// Axons connecting to this neuron
	for index := range neuron.weights {
		if _, ok := prevCol[index].Outputs[neuronIndex]; !ok {
			panic("Neuron weight missing paired output from previous column")
		}
		delete(prevCol[index].Outputs, neuronIndex)
	}

	// Axons leaving from this neuron
	if columnIndex != (len(nn) - 1) {
		nextCol := nn[columnIndex+1]
		for index := range neuron.Outputs {
			delete(nextCol[index].weights, neuronIndex)
		}
	}
}

/**
 * Removes a neuron from an index and places a new neuron there in its place.
 * Effectively resetNeuron + addNeuron, if addNeuron took an index.
 */
func (nn_p *PerceptronNetwork) replaceNeuron(columnIndex int, neuronIndex int, nOpt_p *PerceptronGenerationOptions) {

	nOpt := *nOpt_p
	nn := *nn_p

	neuron := nn[columnIndex][neuronIndex]
	prevCol := nn[columnIndex-1]
	nextCol := nn[columnIndex+1]

	// Delete Axons connecting to this neuron
	for index := range neuron.weights {
		delete(prevCol[index].Outputs, neuronIndex)
	}

	// Delete Axons leaving from this neuron
	for index := range neuron.Outputs {
		delete(nextCol[index].weights, neuronIndex)
	}

	nn[columnIndex][neuronIndex] = Perceptron{
		threshold: nOpt.defaultThreshold,
		Outputs:   make(map[int]bool),
		weights:   make(map[int]float64),
	}

	// Create new Axons connecting to this neuron
	axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxonBack(columnIndex, neuronIndex, nOpt.DefaultAxonWeight)
	}
	// Create new Axons leaving from this neuron
	axonCount = rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxon(columnIndex, neuronIndex, nOpt.DefaultAxonWeight)
	}
}

// We could add a removeNeuron function
// which acted like addNeuron and only removed
// from the end of the list. An old, random-removal
// used to exist but was removed due to performance
// concerns and replaced with replaceNeuron.

/**
 * Add a neuron to the end of a column.
 */
func (nn_p *PerceptronNetwork) addNeuron(columnIndex int, nOpt_p *PerceptronGenerationOptions) {

	nn := *nn_p
	nOpt := *nOpt_p

	nn[columnIndex] = append(nn[columnIndex],
		Perceptron{
			threshold: nOpt.defaultThreshold,
			Outputs:   make(map[int]bool),
			weights:   make(map[int]float64),
		})

	// Axons connecting to this neuron
	axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxonBack(columnIndex, len(nn[columnIndex])-1, nOpt.DefaultAxonWeight)
	}
	// Axons leaving from this neuron
	axonCount = rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
	for i := 0; i < axonCount; i++ {
		nn.addAxon(columnIndex, len(nn[columnIndex])-1, nOpt.DefaultAxonWeight)
	}
}

/**
 * Remove the column between our output column and the column two indexes prior.
 */
func (nn_p *PerceptronNetwork) removeColumn(cOpt_p *PerceptronColumnGenerationOptions) *PerceptronNetwork {

	nn := *nn_p
	cOpt := *cOpt_p
	nOpt := *(cOpt.neuronOptions)

	i := len(nn) - 2

	// Replace the column before output with
	// an empty column, disconnecting all
	// neurons from the previous column
	// in the process.
	for j := 0; j < len(nn[i]); j++ {
		nn.resetNeuron(i, j)
	}

	nn = append(nn[:i], nn[i+1:]...)

	i--

	// Add in some random connections from the
	// new final column to the output column.
	for j := 0; j < len(nn[i]); j++ {
		axonCount := rand.Intn(nOpt.maxAxons-nOpt.minAxons) + nOpt.minAxons
		for k := 0; k < axonCount; k++ {
			nn.addAxon(i, j, nOpt.DefaultAxonWeight)
		}
	}

	return &nn
}

/**
 * Add a column between our output column and the column immediately prior.
 */
func (nn_p *PerceptronNetwork) addColumn(cOpt_p *PerceptronColumnGenerationOptions) *PerceptronNetwork {

	nn := *nn_p
	cOpt := *cOpt_p

	i := len(nn) - 1

	outColumn := nn.copyOutColumn()

	for j := 0; j < len(nn[i]); j++ {
		nn.resetNeuron(i, j)
	}
	nn[i] = []Perceptron{}

	// Place the new outColumn after the
	// empty column.
	nn = append(nn, outColumn)

	newSize := rand.Intn(cOpt.MaxSize-cOpt.MinSize) + cOpt.MinSize

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
func (nn_p *PerceptronNetwork) swapAxons(columnIndex int, neuronIndex1 int, neuronIndex2 int) {

	nn := *nn_p

	neuron1 := nn[columnIndex][neuronIndex1]
	neuron2 := nn[columnIndex][neuronIndex2]

	// This list generation could be removed
	// with a data structure that kept track
	// of a key list along with a map, or
	// otherwise provided random-element
	// access to a map.
	axonList1 := KeySet(neuron1.Outputs)
	axonIndex1 := axonList1[rand.Intn(len(axonList1))]
	axon1 := nn[columnIndex+1][axonIndex1]

	axonList2 := KeySet(neuron2.Outputs)
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
	neuron1.Outputs[axonIndex2] = true
	neuron2.Outputs[axonIndex1] = true
	// Clean up old sendpoints
	delete(neuron1.Outputs, axonIndex1)
	delete(neuron2.Outputs, axonIndex2)
}

/**
 * Remove a random Axon which starts from this neuron's index.
 */
func (nn *PerceptronNetwork) removeAxon(columnIndex int, neuronIndex int) {

	neuron := (*nn)[columnIndex][neuronIndex]

	axonIndex := KeySet(neuron.Outputs)[rand.Intn(len(neuron.Outputs))]
	axon := (*nn)[columnIndex+1][axonIndex]

	// Delete the endpoint
	delete(axon.weights, neuronIndex)

	// Delete the sendpoint
	delete(neuron.Outputs, axonIndex)
}

/**
 * Add a random Axon which starts from this neuron's index.
 */
func (nn_p *PerceptronNetwork) addAxon(columnIndex int, neuronIndex int, DefaultAxonWeight float64) {

	nn := *nn_p

	nextCol := nn[columnIndex+1]
	neuron := nn[columnIndex][neuronIndex]

	if len(neuron.Outputs) >= len(nextCol) {
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
	// already exist in our neuron's Outputs
	for choice := range neuron.Outputs {
		delete(choices, choice)
	}

	choice := KeySet(choices)[rand.Intn(len(choices))]

	nextCol[choice].weights[neuronIndex] = DefaultAxonWeight
	neuron.Outputs[choice] = true
}

/**
 * Add a random Axon which ends at this neuron's index.
 */
func (nn_p *PerceptronNetwork) addAxonBack(columnIndex int, neuronIndex int, DefaultAxonWeight float64) {

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

	nextCol[choice].Outputs[neuronIndex] = true
	neuron.weights[choice] = DefaultAxonWeight
}

func KeySet(m map[int]bool) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}
