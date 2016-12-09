package gp

import (
	"math"
	"math/rand"
)

type Operator func(*GP, ...*Node) int

func neg(gp *GP, x ...*Node) int {
	return Eval(x[0]) * -1
}

func add(gp *GP, xs ...*Node) int {
	return Eval(xs[0]) + Eval(xs[1])
}

func subtract(gp *GP, xs ...*Node) int {
	return Eval(xs[0]) - Eval(xs[1])
}

func multiply(gp *GP, xs ...*Node) int {
	return Eval(xs[0]) * Eval(xs[1])
}

func divide(gp *GP, xs ...*Node) int {
	a := Eval(xs[0])
	b := Eval(xs[1])
	if b != 0 {
		return a / b
	}
	// We can't really punish the GP
	// for randomly deciding to divide by zero.
	return 0
}

func do2(gp *GP, xs ...*Node) int {
	Eval(xs[0])
	return Eval(xs[1])
}

func do3(gp *GP, xs ...*Node) int {
	Eval(xs[0])
	Eval(xs[1])
	return Eval(xs[2])
}

func pow(gp *GP, xs ...*Node) int {
	return int(math.Pow(float64(Eval(xs[0])), float64(Eval(xs[1]))))
}

func pow2(gp *GP, xs ...*Node) int {
	return int(math.Pow(float64(Eval(xs[0])), 2.0))
}

func pow3(gp *GP, xs ...*Node) int {
	return int(math.Pow(float64(Eval(xs[0])), 3.0))
}

func mod(gp *GP, xs ...*Node) int {
	a := Eval(xs[0])
	b := Eval(xs[1])
	if b != 0 {
		return a % b
	}
	return 0
}

func whilePositive(gp *GP, xs ...*Node) int {
	loops := 0
	var out int
	for Eval(xs[0]) > 0 && loops < 10 {
		out = Eval(xs[1])
		loops++
	}
	return out
}

func neZero(gp *GP, xs ...*Node) int {
	if Eval(xs[0]) != 0 {
		return Eval(xs[1])
	}
	return Eval(xs[2])
}

func isPositive(gp *GP, xs ...*Node) int {
	if Eval(xs[0]) > 0 {
		return Eval(xs[1])
	}
	return Eval(xs[2])
}

func ifRand(gp *GP, xs ...*Node) int {
	l := len(xs)
	return Eval(xs[rand.Intn(l)])
}

func randv(gp *GP, nothing ...*Node) int {
	return rand.Intn(10)
}

func zero(gp *GP, nothing ...*Node) int {
	return 0
}

func one(gp *GP, nothing ...*Node) int {
	return 1
}

func two(gp *GP, nothing ...*Node) int {
	return 2
}

func three(gp *GP, nothing ...*Node) int {
	return 3
}

func four(gp *GP, nothing ...*Node) int {
	return 4
}

func five(gp *GP, nothing ...*Node) int {
	return 5
}

func six(gp *GP, nothing ...*Node) int {
	return 6
}

func seven(gp *GP, nothing ...*Node) int {
	return 7
}

func eight(gp *GP, nothing ...*Node) int {
	return 8
}

func nine(gp *GP, nothing ...*Node) int {
	return 9
}
