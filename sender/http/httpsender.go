// Package http provide an implementation of the sender interface to POST json data to a HTTP RESTFul API.
package http // import "goex/ltser/sender/http"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A Sender send json objects to HTTP RESTFul API.
type Sender struct {
	targetURL string
}

// NewSender returns a new Sender to a given HTTP url.
func NewSender(url string) *Sender {
	httpSender := new(Sender)
	httpSender.targetURL = url

	return httpSender
}

// Send POST json objects to target url.
func (s *Sender) Send(b []byte) error {
	r, err := http.Post(s.targetURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %q", r.Status)
	}

	fmt.Print(".")
	return nil
}
