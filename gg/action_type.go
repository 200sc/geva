package gg

type Action func() float64

type ActionType func(val *float64, strength float64) Action

var (
	BaseActionTypes = []ActionType{
		add,
		set,
		mod,
		div,
		mult,
	}
	BaseActionWeights = []float64{
		0.5,
		0.3,
		0.05,
		0.1,
		0.1,
	}
	BaseActionStrength = []Range{
		{0, 15},
		{0, 30},
		{5, 4},
		{0, 10},
		{0, 5},
	}
)

func add(val *float64, strength float64) Action {
	return func() float64 {
		*val += strength
		return *val
	}
}

func set(val *float64, strength float64) Action {
	return func() float64 {
		*val = strength
		return *val
	}
}

func mod(val *float64, strength float64) Action {
	return func() float64 {
		if strength == 0 {
			*val = 0
		} else {
			for *val > strength {
				*val -= strength
			}
			for *val < strength {
				*val += strength
			}
		}
		return *val
	}
}

func div(val *float64, strength float64) Action {
	return func() float64 {
		if strength == 0 {
			*val = 0
		} else {
			*val /= strength
		}
		return *val
	}
}

func mult(val *float64, strength float64) Action {
	return func() float64 {
		*val *= strength
		return *val
	}
}
