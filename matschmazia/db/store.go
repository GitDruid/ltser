// Package db provide an interface to save and read matschmazia sensors' data.
package db // import "goex/ltser/matschmazia/db"

import "goex/ltser/matschmazia/models"

// A Store save matschmazia sensors' data to a target database.
type Store interface {
	Save(sd models.RawData) error
	Read(m models.Measurement, rStart, rStop, station string) (res []float64, err error)
}
