// Package csvjson provide a *csvjson.Reader that wraps a *csv.Reader and returs json []bytes.
// See also json.Encoder.
package csvjson // import "goex/ltser/csvjson"

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	defHeaderRows   = 1
	defIndentFormat = false
	defIntent       = ""
)

// A Reader convert records read by a csv.Reader into json objects.
type Reader struct {
	HeadersRows  uint
	IndentFormat bool
	Intent       string
	cr           *csv.Reader
	headers      []string
	rowsCount    uint
}

// NewReader returns a new Reader that read from r.
func NewReader(r csv.Reader) *Reader {
	csvconvReder := new(Reader)
	csvconvReder.HeadersRows = defHeaderRows
	csvconvReder.IndentFormat = defIndentFormat
	csvconvReder.Intent = defIntent
	csvconvReder.cr = &r

	return csvconvReder
}

// Read obtains one json object from a record read from r.
func (r *Reader) Read() ([]byte, error) {

	// Read headers.
	for r.rowsCount < r.HeadersRows {
		record, err := r.cr.Read()
		r.rowsCount++
		if err != nil {
			return nil, err
		}

		if r.rowsCount == 1 {
			r.headers = append(r.headers, record...) // Cloning values since r.cr.ReuseRecord could be true.
		}
	}

	// Read data.
	record, err := r.cr.Read()
	r.rowsCount++
	if err != nil {
		return nil, err
	}

	// Trasform data.
	if r.rowsCount == 1 && len(r.headers) == 0 { // If headers are missing, generate default columns' names.
		for i := 0; i < len(record); i++ {
			r.headers = append(r.headers, "column"+strconv.Itoa(i))
		}
	}
	m, err := toMap(r.headers, record)
	if err != nil {
		return nil, fmt.Errorf("malformed row #%v (%s)", r.rowsCount, err)
	}
	jsonBytes, err := doMarshal(r, m)
	if err != nil {
		return nil, fmt.Errorf("malformed row #%v (%s)", r.rowsCount, err)
	}

	return jsonBytes, nil
}

func toMap(k []string, v []string) (map[string]string, error) {
	if len(v) != len(k) {
		return nil, errors.New("keys and values sizes don't match")
	}

	m := make(map[string]string)

	for i := 0; i < len(k); i++ {
		m[k[i]] = v[i]
	}

	return m, nil
}

func doMarshal(r *Reader, v interface{}) ([]byte, error) {
	if r.IndentFormat {
		return json.MarshalIndent(v, "", r.Intent)
	} else {
		return json.Marshal(v)
	}
}
