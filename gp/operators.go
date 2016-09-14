package gp

import (
	"math/rand"
)

var operatorNames map[func(*GP, ...Node)]string = map[func(*GP, ...Node)]string{
	neg:        "neg",
	add:        "add",
	subtract:   "subtract",
	multiply:   "multiply",
	divide:     "divide",
	pow:        "pow",
	mod:        "mod",
	neZero:     "NEZ",
	isPositive: "+?",
	ifRand:     "rand?",
	rand:       "rand",
	zero:       "0",
	one:        "1",
	two:        "2",
	three:      "3",
	four:       "four",
	five:       "five",
	six:        "six",
	seven:      "seven",
	eight:      "eight",
	nine:       "nine",
}

func AddOperatorName(op func(*GP, ...Node) int, name string) {
	operatorNames[op] = name
}

func GetOperatorName(op func(*GP, ...Node) int) string {
	if v, ok := operatorNames[op]; ok {
		return v
	}
	return "N/A"
}

func neg(gp *GP, x ...Node) int {
	return Eval(x[0] * -1)
}

func add(gp *GP, xs ...Node) int {
	return Eval(xs[0]) + Eval(xs[1])
}

func subtract(gp *GP, xs ...Node) int {
	return Eval(xs[0]) - Eval(xs[1])
}

func multiply(gp *GP, xs ...Node) int {
	return Eval(xs[0]) * Eval(xs[1])
}

func divide(gp *GP, xs ...Node) int {
	return Eval(xs[0]) / Eval(xs[1])
}

func pow(gp *GP, xs ...Node) int {
	return Eval(xs[0]) ^ Eval(xs[1])
}

func mod(gp *GP, xs ...Node) int {
	return Eval(xs[0]) % Eval(xs[1])
}

func neZero(gp *GP, xs ...Node) int {
	if Eval(xs[0]) != 0 {
		return Eval(xs[1])
	}
	return Eval(xs[2])
}

func isPositive(gp *GP, xs ...Node) int {
	if Eval(xs[0]) > 0 {
		return Eval(xs[1])
	}
	return Eval(xs[2])
}

func ifRand(gp *GP, xs ...Node) int {
	l := len(gp*GP, xs)
	return Eval(xs[rand.IntN(l)])
}

func rand(gp *GP, nothing ...Node) int {
	return rand.IntN(10)
}

func zero(gp *GP, nothing ...Node) int {
	return 0
}

func one(gp *GP, nothing ...Node) int {
	return 1
}

func two(gp *GP, nothing ...Node) int {
	return 2
}

func three(gp *GP, nothing ...Node) int {
	return 3
}

func four(gp *GP, nothing ...Node) int {
	return 4
}

func five(gp *GP, nothing ...Node) int {
	return 5
}

func six(gp *GP, nothing ...Node) int {
	return 6
}

func seven(gp *GP, nothing ...Node) int {
	return 7
}

func eight(gp *GP, nothing ...Node) int {
	return 8
}

func nine(gp *GP, nothing ...Node) int {
	return 9
}
