package influxdb2

import (
	"errors"
	"goex/ltser/matschmazia/models"
	"goex/ltser/timeseries"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// Result implements db.ObservationsIterator allowing to iterate query results.
type Result struct {
	station      models.Station
	measurement  models.Measurement
	queryResult  *influxdb2.QueryTableResult
	currentValue *timeseries.TimeValue
	currentError error
}

// ErrEndOfRecords occurs at the End Of Records.
var ErrEndOfRecords = errors.New("EOR")

// Next allows to obtain next value in the result.
// It will returns err=ErrEndOfRecords if no more records are available.
func (r *Result) Next() (*timeseries.TimeValue, error) {
	if r.currentError != nil {
		return nil, r.currentError
	}

	returnValue := r.currentValue

	if r.queryResult.Next() {
		r.currentValue = &timeseries.TimeValue{
			Time:  r.queryResult.Record().Time(),
			Value: r.queryResult.Record().Value(),
		}
		r.currentError = r.queryResult.Err()
	} else {
		r.currentValue = nil
		r.currentError = ErrEndOfRecords
	}

	return returnValue, nil
}

// Station returns info about station.
func (r *Result) Station() models.Station {
	return r.station
}

// Measurement returns info about measurement.
func (r *Result) Measurement() models.Measurement {
	return r.measurement
}
