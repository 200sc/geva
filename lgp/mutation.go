package lgp

import (
	"math/rand"
)

func (gp *LGP) SwapMutate() {
	// Find two actions and swap them
	i := rand.Intn(len(gp.Instructions))
	j := rand.Intn(len(gp.Instructions))
	if j == i {
		j = (j + 1) % len(gp.Instructions)
	}
	gp.Instructions[i], gp.Instructions[j] = gp.Instructions[j], gp.Instructions[i]
}

func (gp *LGP) ShrinkMutate() {
	if len(gp.Instructions) <= gpOptions.MinActionCount {
		return
	}
	// Get rid of a random action
	i := rand.Intn(len(gp.Instructions))
	gp.Instructions = append(gp.Instructions[:i], gp.Instructions[i+1:]...)
}

func (gp *LGP) ExpandMutate() {
	if len(gp.Instructions) > gpOptions.MaxActionCount {
		return
	}
	// Add an action at a random spot
	inst := gp.GetInstruction()
	i := rand.Intn(len(gp.Instructions))
	half := append(gp.Instructions[:i], inst)
	gp.Instructions = append(half, gp.Instructions[i:]...)
}

func (gp *LGP) ValueMutate() {
	// Mutate an argument set
	i := rand.Intn(len(gp.Instructions))
	inst := gp.Instructions[i]

	gp.Instructions[i].Args = getArgs(inst.Act.Args, len(*gp.Mem)+SPECIAL_REGISTERS)

}
func (gp *LGP) MemMutate() {
	// Mutate a value in memory's start value
	i := rand.Intn(len(*gp.MemStart))
	*(*gp.MemStart)[i] = rand.Intn(10+SPECIAL_REGISTERS) - SPECIAL_REGISTERS
}
