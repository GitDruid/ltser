package timeseries // import "goex/ltser/timeseries"

import (
	"time"
)

// TimeValue represents a single value in time.
type TimeValue struct {
	Time  time.Time
	Value interface{}
}

// TimeSeries represents an historical sequence of values.
type TimeSeries struct {
	Times  []time.Time
	Values []float64
}

// New returns a new TimeSeries.
func New() *TimeSeries {
	s := new(TimeSeries)
	s.Times = make([]time.Time, 0)
	s.Values = make([]float64, 0)

	return s
}

// Lenght returns the lenght of the series.
func (ts *TimeSeries) Lenght() int {
	return len(ts.Values)
}

// AddTimeValue records an observation at a specified time.
func (ts *TimeSeries) AddTimeValue(tv *TimeValue) {
	ts.Times = append(ts.Times, tv.Time)
	ts.Values = append(ts.Values, tv.Value.(float64))
}

// AddWithTime records an observation at a specified time.
func (ts *TimeSeries) AddWithTime(v float64, t time.Time) {
	ts.Times = append(ts.Times, t)
	ts.Values = append(ts.Values, v)
}

// Add records an observation at the current time.
func (ts *TimeSeries) Add(v float64) {
	ts.AddWithTime(v, time.Now())
}
