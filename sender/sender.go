// Package sender provide an interface to send json data to a target.
package sender // import "goex/ltser/sender"

// A Sender send json objects to a target.
type Sender interface {
	Send(b []byte) error
}
