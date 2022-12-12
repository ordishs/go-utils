package utils

import "encoding/hex"

// ReverseSlice reverses the order of the items in the slice.
// The slice is modified in place.
// The slice is returned for convenience.
// The slice can be of any type.
// The slice can be nil.
// The slice can be empty.
// The slice can contain any number of items.
func ReverseSlice[T any](a []T) []T {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// DecideAndReverseHexString decodes the given hex string and then reverses the bytes.
// This is useful for converting Bitcoin hex strings to byte slices in little endian format.
func DecodeAndReverseHexString(hexStr string) ([]byte, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return ReverseSlice(b), nil

}
