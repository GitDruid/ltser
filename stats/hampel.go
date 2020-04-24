package stats

import "math"

// Hampel filter.
// https://towardsdatascience.com/outlier-detection-with-hampel-filter-85ddf523c73d
//
// wSize: size of the sliding window.
// nSigmas: number of standard deviations which identify the outlier.
//
// Select these two parameters depending on the use-case. A higher standard deviation
// threshold makes the filter more forgiving, a lower one identifies more points as outliers.
// Setting the threshold to 0 corresponds to John Tukey’s median filter.
func Hampel(series []float64, wSize, nSigmas int) (newSeries []float64, indexes []int, err error) {
	// For the MAD to be a consistent estimator for the standard deviation,
	// it must be multiplied by a constant scale factor k.
	// The factor is dependent on the distribution, for Gaussian it is approximately 1.4826.
	k := 1.4826
	n := len(series)
	newSeries = copySlice(series)

	for i := wSize; i < n-wSize; i++ {
		mad, median, err := MMAD(series[(i - wSize):(i + wSize)])
		if err != nil {
			return nil, nil, err
		}

		S0 := k * mad

		// If the considered observation differs from the window median by more than x
		// standard deviations, it must be treaten as an outlier and replaced with the median.
		if math.Abs(series[i]-median) > float64(nSigmas)*S0 {
			newSeries[i] = median
			indexes = append(indexes, i)
		}
	}

	return newSeries, indexes, err
}

/*
def hampel_filter_forloop(input_series, window_size, n_sigmas=3):

    n = len(input_series)
    new_series = input_series.copy()
    k = 1.4826 # scale factor for Gaussian distribution

    indices = []

    # possibly use np.nanmedian
    for i in range((window_size),(n - window_size)):
        x0 = np.median(input_series[(i - window_size):(i + window_size)])
        S0 = k * np.median(np.abs(input_series[(i - window_size):(i + window_size)] - x0))
        if (np.abs(input_series[i] - x0) > n_sigmas * S0):
            new_series[i] = x0
            indices.append(i)

    return new_series, indices
*/
