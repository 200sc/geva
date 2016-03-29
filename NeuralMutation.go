package neural

import (
	"rand"
)

type NetworkMutationOptions struct {
	neuronOptions *NeuronMutationOptions
	// per column
	nodeRemovalChance float64 
	nodeAdditionChance float64 
	axonRemovalChance float64 
	axonAdditionChance float64 
	axonSwapChance float64
	// per network
	columnRemovalChance float64 
	columnAdditionChance float64 
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

func mutateFloat(toMutate float64, opt FloatMutationOptions) float {
	if rand.Float64() >= opt.mutChance {
		return toMutate
	}
	
	out := (opt.mutMagnitude * rand.Intn(opt.mutRange*2)) + toMutate
	return out - (opt.mutMagnitude * opt.mutRange)
}

func (n *Neuron) mutate(mOpt_p *NeuronMutationOptions) Neuron {
	mOpt := *mOpt_p

	newThreshold := mutateFloat(n.threshold, mOpt.thresholdOptions)

	newWeights := []float
	for _,weight := range n.weights {
		newWeights = append(weights, mutateFloat(weight, mOpt.weightOptions))
	}

	return Neuron{n.outputs, newThreshold, newWeights}
} 

func (nn *Network) mutate(mOpt_p *NetworkMutationOptions) *Network {
	
	mOpt := *mOpt_p

	newNetwork := nn.copy()

	if rand.Float64() < mOpt.columnRemovalChance {

	}

	if rand.Float64() < mOpt.columnAdditionChance {

	}

	for i := range len(newNetwork.columns-1) {

		col := newNetwork.columns[i]

		if rand.Float64() < mOpt.axonSwapChance {
			neuron1 := col[rand.Intn(len(col))]
			axonIndex1 := neuron.outputs[rand.Intn(len(neuron.ouputs))]
			axonEnd1 := newNetwork.columns[i+1][axonIndex]

			neuron2

		}
		
		if rand.Float64() < mOpt.axonRemovalChance {
			neuronIndex := rand.Intn(len(col))
			neuron := col[neuronIndex]

			// We don't want to remove the only
			// axon a neuron has. 
			if len(neuron.outputs) != 1 {
			
				axonIndex := neuron.outputs[rand.Intn(len(neuron.ouputs))]
				axonEnd := newNetwork.columns[i+1][axonIndex]

				// Delete the endpoint
				delete(axonEnd.weights, neuronIndex)

				// Delete the sendpoint 
				for i,val := range neuron.outputs {
					if val == axonIndex {
						neuron.outputs = neuron.outputs[:i] + neuron.outputs[i+1:]
					}
				}
			}	 
		}

		// Add an additional axon from this column
		// to the next column
		if rand.Float64() < mOpt.axonAdditionChance {
			neuronIndex := rand.Intn(len(col))
			neuron := col[neuronIndex]
			nextCol := newNetwork.columns[i+1]
			
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
			if len(nextCol) != len(neuron.outputs) {

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
				for _,choice := range neuron.outputs {
					choices[choice] = false
				}

				choiceIndex := rand.Intn(len(choices))

				// Our choice is the ith value in choices
				// that is true. 
				for k,v := range choices {
					if v {
						if choiceIndex == 0 {
							// MAGIC NUMBER ALERT
							// DEFAULT WEIGHT IS 1 BUT THAT ISN'T EXPLAINED
							nextCol[k].weights[neuronIndex] = 1
							neuron.outputs = append(neuron.outputs,k)
							break
						}
						choiceIndex--
					}
				}
			}
		}
	}

	// Mutate individual neurons
	for x,column := range nn.columns {
		for y, neuron := range column {
			newNetwork.columns[x][y] = neuron.mutate(mOpt.neuronOptions)
		}
	}

	return newNetwork
}