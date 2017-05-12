package eda

import "bitbucket.org/StephenPatrick/goevo/env"

// A FullBivariateEnv represents a 2d array of Bivariate environments
// where each potential pairing of (a|b) is represented at index [a][b]
type FullBivariateEnv []*BivariateEnv

// NewFullBSBivariateEnv returns a FullBivariateEnv under the assumption
// that the input samples are bitstrings.
func NewFullBSBivariateEnv(samples []*env.F) FullBivariateEnv {
	length := len(*samples[0])
	fbs := make([]*BivariateEnv, length)
	for i := 0; i < length; i++ {
		fbs[i] = NewBSBivariateEnv(samples, i)
	}
	return fbs
}

// BivariateEnv represents the relationship (a|b) for all b
// in some sample set
type BivariateEnv struct {
	domain []float64
	bf     []*env.F
}

// NewBSBivariateEnv returns a bivariate environment
// from samples and a index a for (a|b)
func NewBSBivariateEnv(samples []*env.F, a int) *BivariateEnv {
	be := new(BivariateEnv)
	be.domain = []float64{0.0, 1.0}
	be.bf = make([]*env.F, 2)
	for i := 0; i < 2; i++ {
		be.bf[i] = env.NewF(len(*(samples[0])), 0.0)
	}
	for i := 0; i < len(*(samples[0])); i++ {
		if i == a {
			// The chance that n is true if n is true is 1
			*(*be.bf[0])[a] = 1.0
			// The chance that n is true if n is false is 0
			*(*be.bf[1])[a] = 0.0
		}
		ptt, ptf := BitStringBivariate(samples, a, i)
		*(*be.bf[0])[i] = ptt
		*(*be.bf[1])[i] = ptf
	}
	return be
}

// BitStringBivariate returns the probabilities p(a|b=t) and p(a|b=f)
func BitStringBivariate(samples []*env.F, a, b int) (float64, float64) {
	ptt := 0.0
	ptf := 0.0
	pft := 0.0
	pff := 0.0
	for _, s := range samples {
		af := *(*s)[a]
		bf := *(*s)[b]
		if af == 1.0 {
			if bf == 0.0 {
				ptf++
			} else {
				ptt++
			}
		} else {
			if bf == 0.0 {
				pff++
			} else {
				pft++
			}
		}
	}
	ptt /= float64(len(samples))
	ptf /= float64(len(samples))
	pft /= float64(len(samples))
	pff /= float64(len(samples))
	//fmt.Println(ptt, ptf, pft, pff)
	//fmt.Println(ptt/(ptt+pft), ptf/(ptf+pff))
	if ptt+pft > 0 {
		ptt = ptt / (ptt + pft)
	} // else ptt is already 0
	if ptf+pff > 0 {
		ptf = ptf / (ptf + pff)
	} // else ptf is already 0
	//fmt.Println(ptt, ptf)
	return ptt, ptf
}
