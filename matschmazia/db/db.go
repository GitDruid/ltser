// Package db provide interfaces to save and read matschmazia sensors' data.
package db // import "goex/ltser/matschmazia/db"

import (
	"goex/ltser/matschmazia/models"
	"time"
)

// A Writer save matschmazia sensors' data.
type Writer interface {
	Write(sd models.RawData) error
}

// Iterator allows to iterate a sequence of float64 values.
type Iterator interface {
	Next() (n float64, err error)
}

// A Reader read matschmazia sensors' data.
type Reader interface {
	Read(m models.Measurement, rStart, rStop time.Time, station string) (Iterator, error)
	ReadAll(m models.Measurement, rStart, rStop time.Time, station string) ([]float64, error)
}

// A ReadWriter read and save matschmazia sensors' data.
type ReadWriter interface {
	Reader
	Writer
}
