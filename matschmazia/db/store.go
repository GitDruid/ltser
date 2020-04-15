// Package db provide an interface to save matschmazia sensors' data to database.
package db // import "goex/ltser/matschmazia/db"

import "goex/ltser/matschmazia/models"

// A Store save matschmazia sensors' data to a target database.
type Store interface {
	Save(sd models.SensorData) error
}
