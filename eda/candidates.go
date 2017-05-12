package eda

import (
	"bitbucket.org/StephenPatrick/goevo/env"
)

// BestCandidates represents the top N candidates sampled
// from an EDA.
// Todo: This should be a heap, not a linked list,
// for sufficiently large sizes
type BestCandidates struct {
	Front  *Candidate
	Back   *Candidate
	Length int
	Limit  int
}

// Candidate is an individual candidate in a BestCandidates list
type Candidate struct {
	*env.F
	Fitness int
	Next    *Candidate
	Prev    *Candidate
}

// NewBestCandidates creates a default BestCandidates
// with model.samples samples added to the candidate list.
// if sFunc is not supplied, GetSample on a copy of the model's
// initial environment will be used.
func NewBestCandidates(m Model, bcsLimit int, sFunc func() *env.F) *BestCandidates {
	bm := m.BaseModel()

	bcs := new(BestCandidates)
	bcs.Limit = bcsLimit
	eCopy := bm.F.Copy()
	for i := 0; i < bm.samples; i++ {
		// We set the sample to cga.F right now
		// as our fitness function takes in a model
		// this might change
		if sFunc == nil {
			bm.F = GetSample(eCopy)
		} else {
			bm.F = sFunc()
		}
		bcs.Add(bm.fitness(bm), bm.F)
	}

	bm.F = eCopy
	return bcs
}

// Slice converts BestCandidates from a linked list to a slice
func (bc *BestCandidates) Slice() []*env.F {
	sl := make([]*env.F, bc.Length)
	i := 0
	for cur := bc.Front; cur != nil; cur = cur.Next {
		sl[i] = cur.F
		i++
	}
	return sl
}

// Add appends a new candidate to the best candidates,
// if it is better than any existing candidate
func (bc *BestCandidates) Add(f int, c *env.F) {
	if bc.Front == nil {
		bc.Front = &Candidate{c, f, nil, nil}
		bc.Back = bc.Front
		bc.Length++
	} else if bc.Limit > bc.Length {
		bc.add(&Candidate{c, f, nil, nil})
		bc.Length++
	} else if bc.Back.Fitness > f {
		bc.add(&Candidate{c, f, nil, nil})
		bc.Back = bc.Back.Prev
		bc.Back.Next = nil
	}
	// else do not add this
}

func (bc *BestCandidates) add(cand *Candidate) {
	cur := bc.Front
	for cur != nil && cur.Fitness < cand.Fitness {
		cur = cur.Next
	}
	if cur == nil {
		bc.Back.Next = cand
		cand.Prev = bc.Back
		bc.Back = cand
	} else {
		if cur.Prev != nil {
			cur.Prev.Next = cand
		}
		cand.Prev = cur.Prev
		cur.Prev = cand
		cand.Next = cur
		if cur == bc.Front {
			bc.Front = cand
		}
	}
}
