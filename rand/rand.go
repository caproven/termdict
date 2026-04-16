package rand

import (
	"math/rand/v2"
)

type Default struct{}

func (Default) IntN(n int) int {
	return rand.IntN(n)
}

type Rand struct {
	base *rand.Rand
}

func NewRand(seed uint64) *Rand {
	return &Rand{
		base: rand.New(rand.NewPCG(seed, seed)),
	}
}

func (r *Rand) IntN(n int) int {
	return r.base.IntN(n)
}
