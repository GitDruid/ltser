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
	defer r.Body.Close()

	// TODO: Not safe implementation.
	// See: https://haisum.github.io/2017/09/11/golang-ioutil-readall/
	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %q", r.Status)
	}

	return nil
}
