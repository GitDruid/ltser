package timeseries

import "fmt"

// Float attaches the methods of Observable to a float64.
type Float float64

// NewFloat returns a Float.
func NewFloat() Observable {
	f := Float(0)
	return &f
}

// String returns the float as a string.
func (f *Float) String() string {
	return fmt.Sprintf("%g", f.Value())
}

// Value returns the float's value.
func (f *Float) Value() float64 {
	return float64(*f)
}

// Multiply a Float for a given ratio.
func (f *Float) Multiply(ratio float64) {
	*f *= Float(ratio)
}

// Add another Float's value to the current one.
func (f *Float) Add(other Observable) {
	o := other.(*Float)
	*f += *o
}

// Clear the observation so it can be reused.
func (f *Float) Clear() {
	*f = 0
}

// CopyFrom copies the contents of a given observation to self.
func (f *Float) CopyFrom(other Observable) {
	o := other.(*Float)
	*f = *o
}
