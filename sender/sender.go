// Package sender provide an interface to send json data to a target.
//
// TODO: add a flag to switch between sequential implementation (that guarantees
// rows order) and parallel implementation (no order guaranteed).
// Move the implementation in a sender package with:
// 		func Send(b []byte)
// 		func WaitForCompletion()
//		func State() []error
//		MaxConcurrency uint
//		TargetURL string
package sender // import "goex/ltser/sender"

// A Sender send json objects to a target.
type Sender interface {
	Send(b []byte) error
}
