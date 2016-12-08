package gp

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
	v := int(math.Pow(float64(gp.getVal(xs[1])), float64(gp.getVal(xs[2]))))
	gp.setReg(xs[0], v)
}

func pow2(gp *LGP, xs ...int) {
	return int(math.Pow(float64(Eval(xs[0])), 2.0))
}

func pow3(gp *LGP, xs ...int) {
	return int(math.Pow(float64(Eval(xs[0])), 3.0))
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
	gp.pc = gp.regVal(xs[0])
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

func (gp *LGP) regVal(r1 int) int {
	return gp.Mem[gp.getReg(r1)]
}

func (gp *LGP) getReg(r1 int) (r2 int) {
	r2 = r1
	// Special coded pointers evaluated here
	if r2 == -1 {
		r2 = gp.lastRegister
	}
	return
}

func (gp *LGP) setReg(r1, v int) {
	gp.lastRegister = gp.getReg(r1)
	gp.Mem[gp.lastRegister] = v
}
