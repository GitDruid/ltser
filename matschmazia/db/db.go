// Package db provide interfaces to save and read matschmazia sensors' data.
package db // import "goex/ltser/matschmazia/db"

import (
	"goex/ltser/matschmazia/models"
	"goex/ltser/timeseries"
	"time"
)

// A Writer save matschmazia sensors' data.
type Writer interface {
	Write(sd models.RawData) error
	WriteObservations(o models.Observations) error
}

// ObservationsIterator allows to iterate a sequence of TimeValues.
type ObservationsIterator interface {
	Next() (*timeseries.TimeValue, error)
	Station() models.Station
	Measurement() models.Measurement
}

// A Reader read matschmazia sensors' data.
type Reader interface {
	Read(m models.Measurement, rStart, rStop time.Time, station string) (ObservationsIterator, error)
	ReadAll(m models.Measurement, rStart, rStop time.Time, station string) (*models.Observations, error)
}

// A ReadWriter read and save matschmazia sensors' data.
type ReadWriter interface {
	Reader
	Writer
}
