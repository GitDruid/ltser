package stats

import (
	"math"
	"sort"
)

// Median gets the median number in a slice of numbers.
// Extracted from https://github.com/montanaflynn/stats
func Median(input []float64) (median float64, err error) {

	// No math is needed if there are no numbers.
	l := len(input)
	if l == 0 {
		return math.NaN(), ErrEmptyInput
	}

	// Start by sorting a copy of the slice.
	c := sortedCopy(input)
	middle := l / 2

	// For even numbers we add the two middle numbers and divide by two.
	// For odd numbers we just use the middle number.
	if l%2 == 0 {
		median = (c[middle-1] + c[middle]) / 2
	} else {
		median = c[middle]
	}

	return median, nil
}

// copySlice copies a slice of float64s
func copySlice(input []float64) []float64 {
	s := make([]float64, len(input))
	copy(s, input)
	return s
}

// sortedCopy returns a sorted copy of float64s
func sortedCopy(input []float64) (copy []float64) {
	copy = copySlice(input)
	sort.Float64s(copy)
	return
}
