package extensions // import "goex/ltser/extensions"

import (
	"time"
)

// DateTimeFlag implements flag.Value interface for parameters in time.RFC3339 format.
type DateTimeFlag struct {
	value time.Time
}

// String method of flag.Value interface.
func (n *DateTimeFlag) String() string {
	return n.value.Format(time.RFC3339)
}

// Set method of flag.Value interface.
func (n *DateTimeFlag) Set(value string) (err error) {
	n.value, err = time.Parse(time.RFC3339, value)
	return err
}

// Value returns the embedded time.Time value.
func (n *DateTimeFlag) Value() time.Time {
	return n.value
}
