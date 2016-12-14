package lgp

import (
	"math"
	"math/rand"
)

type Operator func(*LGP, ...int)

func neg(gp *LGP, xs ...int) {
	gp.setReg(xs[0], gp.regVal(xs[0])*-1)
}

func add(gp *LGP, xs ...int) {
	gp.setReg(xs[0], gp.regVal(xs[1])+gp.regVal(xs[2]))
}

func subtract(gp *LGP, xs ...int) {
	gp.setReg(xs[0], gp.regVal(xs[1])-gp.regVal(xs[2]))
}

func multiply(gp *LGP, xs ...int) {
	gp.setReg(xs[0], gp.regVal(xs[1])*gp.regVal(xs[2]))
}

func divide(gp *LGP, xs ...int) {
	a := gp.regVal(xs[1])
	b := gp.regVal(xs[2])
	if b != 0 {
		gp.setReg(xs[0], a/b)
	}
}

func pow(gp *LGP, xs ...int) {
	v := int(math.Pow(float64(gp.regVal(xs[1])), float64(gp.regVal(xs[2]))))
	gp.setReg(xs[0], v)
}

func pow2(gp *LGP, xs ...int) {
	v := int(math.Pow(float64(gp.regVal(xs[1])), 2))
	gp.setReg(xs[0], v)
}

func pow3(gp *LGP, xs ...int) {
	v := int(math.Pow(float64(gp.regVal(xs[1])), 3))
	gp.setReg(xs[0], v)
}

func mod(gp *LGP, xs ...int) {
	a := gp.regVal(xs[1])
	b := gp.regVal(xs[2])
	if b != 0 {
		gp.setReg(xs[0], a%b)
	}
}

// Branch functions skip the following instruction if true.
func bnez(gp *LGP, xs ...int) {
	if gp.getReg(xs[0]) != 0 {
		gp.pc++
	}
}

func bgz(gp *LGP, xs ...int) {
	if gp.getReg(xs[0]) > 0 {
		gp.pc++
	}
}

func jmp(gp *LGP, xs ...int) {
	pc := gp.regVal(xs[0])
	if pc >= len(gp.Instructions) {
		pc = len(gp.Instructions) - 1
	}
	if pc < 0 {
		pc = 0
	}
	gp.pc = pc
}

func randv(gp *LGP, xs ...int) {
	gp.setReg(xs[0], rand.Intn(10))
}

func zero(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 0)
}

func one(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 1)
}

func two(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 2)
}

func three(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 3)
}

func four(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 4)
}

func five(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 5)
}

func six(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 6)
}

func seven(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 7)
}

func eight(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 8)
}

func nine(gp *LGP, xs ...int) {
	gp.setReg(xs[0], 9)
}

func getEnv(gp *LGP, xs ...int) {
	index := gp.regVal(xs[1])
	if index >= len(*gp.Env) {
		index = len(*gp.Env) - 1
	}
	if index < 0 {
		index = 0
	}
	gp.setReg(xs[0], *(*gp.Env)[index])
}

func setEnv(gp *LGP, xs ...int) {
	index := gp.regVal(xs[1])
	if index >= len(*gp.Env) {
		index = len(*gp.Env) - 1
	}
	if index < 0 {
		index = 0
	}
	*(*gp.Env)[index] = gp.regVal(xs[0])
}

func envLen(gp *LGP, xs ...int) {
	gp.setReg(xs[0], len(*gp.Env))
}

func (gp *LGP) regVal(r1 int) int {
	return *(*gp.Mem)[gp.getReg(r1)]
}

func (gp *LGP) getReg(r1 int) (r2 int) {
	r2 = r1
	// Special coded pointers evaluated here
	if r2 == LAST_WRITTEN {
		r2 = gp.lastRegister
	}
	return
}

func (gp *LGP) setReg(r1, v int) {
	gp.lastRegister = gp.getReg(r1)
	*(*gp.Mem)[gp.lastRegister] = v
}
