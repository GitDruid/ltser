package stats

import "math"

// AbsoluteDeviation calculates absolute difference between each element of a data set and
// a given point (typically a central value like the median or the mean of the data set).
func AbsoluteDeviation(series []float64, x float64) []float64 {
	s := make([]float64, len(series))
	for i, v := range series {
		s[i] = math.Abs(v - x)
	}
	return s
}
