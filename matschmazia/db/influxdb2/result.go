package influxdb2

import (
	"errors"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// Result contains the query result in an iterable form.
type Result struct {
	queryResult *influxdb2.QueryTableResult
}

// ErrEndOfRecords occurs at the End Of Records.
var ErrEndOfRecords = errors.New("EOR")

// Next allows to obtain next value in the result.
// It will returns err=ErrEndOfRecords if no more records are available.
func (r Result) Next() (n float64, err error) {
	if r.queryResult.Next() {
		// Observe when there is new grouping key producing new table
		if r.queryResult.TableChanged() {
			fmt.Printf("table: %s\n", r.queryResult.TableMetadata().String())
		}
		n = r.queryResult.Record().Value().(float64)
		fmt.Printf("row: %s\n", r.queryResult.Record())
		err = r.queryResult.Err()
	} else {
		err = ErrEndOfRecords
	}

	return
}
