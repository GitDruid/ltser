// Package sender provide an interface to send json data to a target.
//
// TODO: add a MaxConcurrency field (default 1) to switch between sequential implementation
// (that guarantees order) and parallel implementation (no order guaranteed).
// Add following methods:
// 		func Send(b []byte)			starts a goroutine to send data
// 		func WaitForCompletion()	waits for all running goroutines to complete
//		func State() []error		returns a list of encountered errors
//		MaxConcurrency uint			specifies how many goroutine can be called concurrently
package sender // import "goex/ltser/sender"

// A Sender send json objects to a target.
type Sender interface {
	Send(b []byte) error
}
