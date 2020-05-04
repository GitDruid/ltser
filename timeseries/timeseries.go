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
type TimeSeries []TimeValue

// AddWithTime records an observation at a specified time.
func (ts *TimeSeries) AddWithTime(tv TimeValue) {
	*ts = append(*ts, tv)
}

// Add records an observation at the current time.
func (ts *TimeSeries) Add(v interface{}) {
	tv := TimeValue{Time: time.Now(), Value: v}
	ts.AddWithTime(tv)
}

// Values returns the raw values of the series.
// TODO: temporary implementation.
func (ts *TimeSeries) Values() []interface{} {
	res := make([]interface{}, len([]TimeValue(*ts)))
	for i, v := range *ts {
		res[i] = v.Value
	}
	return res
}

// FloatValues returns the float64 values of the series.
// TODO: temporary implementation.
func (ts *TimeSeries) FloatValues() []float64 {
	res := make([]float64, len([]TimeValue(*ts)))
	for i, v := range *ts {
		res[i] = v.Value.(float64)
	}
	return res
}
