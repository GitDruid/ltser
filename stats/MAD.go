package stats

import "math"

// MMAD calculate both the Median and the Median Absolute Deviation (around the median) of a series.
// https://en.wikipedia.org/wiki/Median_absolute_deviation
func MMAD(series []float64) (mad float64, median float64, err error) {
	median, err = Median(series)
	if err != nil {
		return math.NaN(), math.NaN(), err
	}
	mad, err = MAD(series, median)
	if err != nil {
		return math.NaN(), math.NaN(), err
	}
	return
}

// MAD calculate the Median Absolute Deviation of a series around a given
// central point (e.g. an already calculated Median or Mean).
func MAD(series []float64, central float64) (mad float64, err error) {
	mad, err = Median(AbsoluteDeviation(series, central))
	if err != nil {
		return math.NaN(), err
	}
	return
}
