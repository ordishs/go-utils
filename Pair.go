package utils

import (
	"fmt"
)

// Pair is a simple struct that holds two values.
type Pair[F, S any] struct {
	First  F
	Second S
}

func NewPair(f, s any) Pair[any, any] {
	return Pair[any, any]{f, s}
}

func (p Pair[F, S]) String() string {
	return fmt.Sprintf("[%v, %v]", p.First, p.Second)
}
