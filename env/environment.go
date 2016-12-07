package env

type Env interface {
	Copy() Env
	Diff([]float64) Env
	New([]float64) Env
}
