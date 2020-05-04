package models // import "goex/ltser/matschmazia/models"

import "goex/ltser/timeseries"

// Measurement represents a measure type.
type Measurement struct {
	name     string
	unit     string
	interval string
}

// Unit is the measurement unit.
func (m Measurement) Unit() string {
	return m.unit
}

// Name is the measurement name.
func (m Measurement) Name() string {
	return m.name
}

func (m Measurement) String() string {
	return m.name
}

// Available measurement.
var (
	Temperature    = Measurement{"temperature", "Celsius", "15 min average"}
	WindSpeed      = Measurement{"wind_speed", "m/s", ""} // Undocumented interval.
	WindGust       = Measurement{"wind_gust", "m/s", ""}  // Undocumented interval.
	Humidity       = Measurement{"humidity", "Percentage", "15 min average"}
	Precipitations = Measurement{"precipitations", "mm", "15 min cumulative sum"}
	Snow           = Measurement{"snow", "m", ""} // Undocumented interval.
)

// Location represents a geographic position.
type Location struct {
	Altitude  int
	Latitude  float32
	Longitude float32
}

// Station represents a sensor station.
type Station struct {
	Location
	Name string
}

// Observations represents a collection of historical Measures from a given Matsch/Mazia sensor station.
type Observations struct {
	Station
	Measurement
	Measures timeseries.TimeSeries
}
