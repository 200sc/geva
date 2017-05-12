package eda

// PBIL is an EDA built on the Population based incremental learning algorithm
type PBIL struct {
	Base
}

// Adjust for a PBIL takes the learningSamples best samples from the pbil's
// distribution and reinforces the distribution in the direction of each of the
// taken samples by learningRate/learningSamples
func (pbil *PBIL) Adjust() Model {
	bcs := NewBestCandidates(pbil, pbil.learningSamples, nil)
	pbil.F.Reinforce(pbil.learningRate/float64(bcs.Length), bcs.Slice()...)
	return pbil
}

// PBILModel initializes a PBIL EDA
func PBILModel(opts ...Option) (Model, error) {
	var err error
	pbil := new(PBIL)
	pbil.Base, err = DefaultBase(opts...)
	return pbil, err
}
