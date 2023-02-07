package utils

import (
	"testing"
	"time"
)

func TestMedianTime(t *testing.T) {
	// Test the median time calculation.
	expected := time.Unix(1675721369, 0)

	timestamps := []int{
		1675725269, // 778028
		1675725209, // 778027
		1675724249, // 778026
		1675723109, // 778025
		1675721789, // 778024
		1675721369, // 778023
		1675720529, // 778022
		1675719149, // 778021
		1675718489, // 778020
		1675717709, // 778019
		1675716929, // 778018
	}

	mt, err := CalcPastMedianTime(timestamps)
	if err != nil {
		t.Errorf("CalcPastMedianTime: %v", err)
	}

	if !mt.Equal(expected) {
		t.Errorf("CalcPastMedianTime: unexpected result - got %v, want %v", mt, expected)
	}
}

func TestBlock5MedianTime(t *testing.T) {
	// Test the median time calculation.
	expected := time.Unix(1231470173, 0)

	timestamps := []int{
		1231471428, // 5
		1231470988, // 4
		1231470173, // 3
		1231469744, // 2
		1231469665, // 1
		1231006505, // Genesis
	}

	mt, err := CalcPastMedianTime(timestamps)
	if err != nil {
		t.Errorf("CalcPastMedianTime: %v", err)
	}

	if !mt.Equal(expected) {
		t.Errorf("CalcPastMedianTime: unexpected result - got %v, want %v", mt, expected)
	}
}
