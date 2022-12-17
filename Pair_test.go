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
