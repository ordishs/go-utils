package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseSlice(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := ReverseSlice(a)

	assert.Equal(t, []int{5, 4, 3, 2, 1}, b)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}
func TestReverseHash(t *testing.T) {
	var hash1 [32]byte
	var hash2 [32]byte
	var reverseHash [32]byte

	for i := 0; i < 32; i++ {
		hash1[i] = byte(i)
		hash2[i] = byte(i)
		reverseHash[31-i] = byte(i)
	}

	reversed := ReverseHash(hash1)

	assert.Equal(t, reverseHash, reversed)
	assert.Equal(t, hash2, hash1)
}

func TestReverseSliceInPlace(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	ReverseSliceInPlace(a)

	assert.Equal(t, []int{5, 4, 3, 2, 1}, a)
}

func TestReverseHashInPlace(t *testing.T) {
	var hash [32]byte
	var reverseHash [32]byte

	for i := 0; i < 32; i++ {
		hash[i] = byte(i)
		reverseHash[31-i] = byte(i)
	}

	ReverseHashInPlace(&hash)

	assert.Equal(t, reverseHash, hash)
}

func TestDecodeAndReverseHexString(t *testing.T) {
	b, err := DecodeAndReverseHexString("0102030405")
	assert.Nil(t, err)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, b)
}

func TestDecodeAndReverseHashString(t *testing.T) {
	hash := make([]byte, 32)
	var reverseHash [32]byte

	for i := 0; i < 32; i++ {
		hash[i] = byte(i)
		reverseHash[31-i] = byte(i)
	}

	hashStr := hex.EncodeToString(hash)

	b, err := DecodeAndReverseHashString(hashStr)
	assert.Nil(t, err)
	assert.Equal(t, reverseHash, b)
}

func TestHexEncodeAndReverseBytes(t *testing.T) {
	str := HexEncodeAndReverseBytes([]byte{5, 4, 3, 2, 1})
	assert.Equal(t, "0102030405", str)
}
