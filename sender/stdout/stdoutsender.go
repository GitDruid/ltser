// Package stdout provide an implementation of the sender interface to write json data to StdOut.
package stdout // import "goex/ltser/sender/stdout"

import (
	"fmt"
	"os"
)

// A Sender send json objects to StdOut.
type Sender struct {
}

// NewSender returns a new Sender to StdOut.
func NewSender() *Sender {
	stdioSender := new(Sender)

	return stdioSender
}

// Send write json objects to StdOut.
func (s *Sender) Send(b []byte) error {
	_, err := fmt.Fprintf(os.Stdout, "%s\n", b)

	return err
}
