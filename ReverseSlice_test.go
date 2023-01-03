package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseSlice(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := ReverseSlice(a)

	assert.Equal(t, []int{5, 4, 3, 2, 1}, b)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestReverseSliceInPlace(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	ReverseSliceInPlace(a)

	assert.Equal(t, []int{5, 4, 3, 2, 1}, a)
}

func TestDecodeAndReverseHexString(t *testing.T) {
	b, err := DecodeAndReverseHexString("0102030405")
	assert.Nil(t, err)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, b)
}

func TestHexEncodeAndReverseBytes(t *testing.T) {
	str := HexEncodeAndReverseBytes([]byte{5, 4, 3, 2, 1})
	assert.Equal(t, "0102030405", str)
}
