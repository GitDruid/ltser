package stats_test

import (
	"goex/ltser/stats"
	"testing"
	//"github.com/montanaflynn/stats"
)

func TestMedian(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{5, 3, 4, 2, 1}, 3.0},
		{[]float64{6, 3, 2, 4, 5, 1}, 3.5},
		{[]float64{1}, 1.0},
		{[]float64{1, 3}, 2.0},
	} {
		got, _ := stats.Median(c.in)
		if got != c.out {
			t.Errorf("Median(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}

	_, err := stats.Median([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}
