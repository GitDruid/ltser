// Package db provide interfaces to save and read matschmazia sensors' data.
package db // import "goex/ltser/matschmazia/db"

import "goex/ltser/matschmazia/models"

// A Writer save matschmazia sensors' data.
type Writer interface {
	Write(sd models.RawData) error
}

// A Reader read matschmazia sensors' data.
type Reader interface {
	Read(m models.Measurement, rStart, rStop, station string) (res []float64, err error)
}

// A ReadWriter read and save matschmazia sensors' data.
type ReadWriter interface {
	Reader
	Writer
}
