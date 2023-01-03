package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair(t *testing.T) {
	p := NewPair("Hello", 123)
	s := fmt.Sprintf("%v", p)
	assert.Equal(t, "[Hello, 123]", s)
}

func TestPair2(t *testing.T) {
	ch := make(chan Pair[string, int])

	pair := NewPair("Hello", 123)

	go func() {
		ch <- pair
	}()

	p := <-ch
	assert.Equal(t, p, pair)

}
