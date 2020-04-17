package extensions // import "goex/ltser/extensions"

import (
	"fmt"
	"strconv"
)

// NotZeroUint32 implements flag.Value interface.
// To be used when some flags values need to be greater than 1.
type NotZeroUint32 struct {
	number uint32
}

// String method of flag.Value interface.
func (n *NotZeroUint32) String() string {
	return fmt.Sprint(n.number)
}

// Set method of flag.Value interface.
func (n *NotZeroUint32) Set(value string) error {
	uintNum, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return err
	}
	if uintNum < 1 {
		return fmt.Errorf("number %v is minor than 1", uintNum)
	}
	n.number = uint32(uintNum)
	return nil
}

// Value returns the embedded uint32 number, using 1 as default value.
func (n *NotZeroUint32) Value() uint32 {
	if n.number > 1 {
		return n.number
	}
	return 1 // Default value.
}
