package env

// An Env represents a slice of values
// which can be modified by a slice of float64s
// or return a patterned difference between
// a slice of float64s and the Env
type Env interface {
	Copy() Env
	Diff([]float64) Env
	MatchDiff([]float64) Env
	New([]float64) Env
}
