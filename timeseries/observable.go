package timeseries

// An Observable is a kind of data that can be aggregated in a time series.
type Observable interface {
	Multiply(ratio float64)    // Multiplies the data in self by a given ratio
	Add(other Observable)      // Adds the data from a different observation to self
	Clear()                    // Clears the observation so it can be reused.
	CopyFrom(other Observable) // Copies the contents of a given observation to self
}

/*
// TimeValue represents a single value in time.
type TimeValue struct {
	Time  time.Time
	Value Observable
}

// TimeSeries represents an historical sequence of values.
type TimeSeries struct {
	times  []time.Time
	values []Observable
}

// New returns a new TimeSeries.
func New() *TimeSeries {
	s := new(TimeSeries)
	s.times = make([]time.Time, 0)
	s.values = make([]Observable, 0)

	return s
}

// AddWithTime records an observation at a specified time.
func (ts *TimeSeries) AddWithTime(tv TimeValue) {
	ts.times = append(ts.times, tv.Time)
	ts.values = append(ts.values, tv.Value)
}

// Add records an observation at the current time.
func (ts *TimeSeries) Add(v Observable) {
	ts.AddWithTime(v, time.Now())
}

// Values returns the raw values of the series.
func (ts *TimeSeries) Values() []Observable {
	return ts.values
}

// FloatValues returns the float64 values of the series.
// TODO: temporary implementation.
func (ts *TimeSeries) FloatValues() []float64 {
	b := make([]float64, len(ts.values))
	for i := range ts.values {
		b[i] = ts.values[i].(*Float).Value()
	}
	return b
}
*/
