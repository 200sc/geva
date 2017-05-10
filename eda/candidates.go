package eda

import "bitbucket.org/StephenPatrick/goevo/env"

type BestCandidates struct {
	Front  *Candidate
	Back   *Candidate
	Length int
	Limit  int
}

type Candidate struct {
	*env.F
	Fitness int
	Next    *Candidate
	Prev    *Candidate
}

func NewBestCandidates(l int) *BestCandidates {
	bc := new(BestCandidates)
	bc.Limit = l
	return bc
}

func (bc *BestCandidates) Slice() []*env.F {
	sl := make([]*env.F, bc.Length)
	i := 0
	for cur := bc.Front; cur != nil; cur = cur.Next {
		sl[i] = cur.F
		i++
	}
	return sl
}

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
