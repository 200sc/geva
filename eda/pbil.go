package eda

type PBIL struct {
	Base
}

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

func PBILModel(opts ...Option) (Model, error) {
	var err error
	pbil := new(PBIL)
	pbil.Base, err = DefaultBase(opts...)
	return pbil, err
}
