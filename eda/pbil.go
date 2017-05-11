package eda

// PBIL is an EDA built on the Population based incremental learning algorithm
type PBIL struct {
	Base
}

// Adjust for a PBIL takes the learningSamples best samples from the pbil's
// distribution and reinforces the distribution in the direction of each of the
// taken samples by learningRate/learningSamples
func (pbil *PBIL) Adjust() Model {

	bcs := NewBestCandidates(pbil.learningSamples)
	eCopy := pbil.F.Copy()
	for i := 0; i < pbil.samples; i++ {
		// We set the sample to pbil.F right now
		// as our fitness function takes in a model
		// this might change
		pbil.F = GetSample(eCopy)
		bcs.Add(pbil.fitness(pbil), pbil.F)
	}
	pbil.F = eCopy
	bcsList := bcs.Slice()
	// Hypothetically bcsList has a length equal to
	// pbil.learningSamples but if samples < learningSamples
	// this this case ensures we still learn a total of learningRate.
	realRate := pbil.learningRate / float64(len(bcsList))
	for _, cand := range bcsList {
		pbil.F.Reinforce(cand, realRate)
	}
	pbil.F.Mutate(pbil.mutationRate, pbil.fmutator)
	pbil.learningRate = pbil.lmutator(pbil.learningRate)
	return pbil
}

// PBILModel initializes a PBIL EDA
func PBILModel(opts ...Option) (Model, error) {
	var err error
	pbil := new(PBIL)
	pbil.Base, err = DefaultBase(opts...)
	return pbil, err
}
