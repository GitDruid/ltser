package models

import "time"

// Measurement represent a measure.
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

// Location represent a geographic position.
type Location struct {
	Altitude  int
	Latitude  float32
	Longitude float32
}

// Station represent a sensor station.
type Station struct {
	Location
	Name string
}

// Measure represent a single measure.
type Measure struct {
	Measurement
	Value float32
}

// SimplePoint represent a Measure in time from a sensor.
type SimplePoint struct {
	Station
	Time time.Time
	Measure
}

// MultiValuePoint represent a point with multiple measures.
type MultiValuePoint struct {
	Station
	Time     time.Time
	Measures []Measure
}
