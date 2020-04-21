// Package influxdb2 provide an implementation of the store interface to save data to database.
package influxdb2 // import "goex/ltser/matschmazia/db/influxdb2"

import (
	"context"
	"goex/ltser/matschmazia/models"
	"math"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// A Store save data to database.
type Store struct {
	url    string
	org    string
	bucket string
	token  string
}

const (
	mTemperature    = "temperature"
	mWind           = "wind"
	mHumidity       = "humidity"
	mPrecipitations = "precipitations"
	mSnow           = "snow"
)

// NewStore returns a new InfluxDB Store.
func NewStore(url, org, bucket, token string) *Store {
	influxDbStore := new(Store)
	influxDbStore.url = url
	influxDbStore.org = org
	influxDbStore.bucket = bucket
	influxDbStore.token = token

	return influxDbStore
}

// Save store data in InfluxDB measurements.
func (s *Store) Save(sd models.SensorData) error {

	client := influxdb2.NewClient(s.url, s.token)
	defer client.Close() // Ensures background processes finishes.

	writeAPI := client.WriteApiBlocking(s.org, s.bucket)

	t, err := time.Parse("2006-01-02 15:04:05 -0700", sd.Time+" +0100") // Measurement time is (UTC +1).
	if err != nil {
		return err
	}

	var points = make([]*influxdb2.Point, 0, 5)

	// Temperature.
	if temp, err := strconv.ParseFloat(sd.AirTempAvg, 32); err == nil && !math.IsNaN(temp) {
		p := influxdb2.NewPointWithMeasurement(mTemperature).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "celsius").
			AddField("avg15", temp).
			SetTime(t)
		points = append(points, p)
	}

	// Wind.
	wsa, err := strconv.ParseFloat(sd.WindSpeedAvg, 32)
	if err != nil || math.IsNaN(wsa) {
		wsa, err = strconv.ParseFloat(sd.WindSpeed, 32) // In same cases values are in sd.WindSpeed, in others in sd.WindSpeedAvg.
	}
	if err == nil && !math.IsNaN(wsa) {
		if wsm, err := strconv.ParseFloat(sd.WindSpeedMax, 32); err == nil && !math.IsNaN(wsm) {
			p := influxdb2.NewPointWithMeasurement(mWind).
				AddTag("station", sd.Station).
				AddTag("altitude", sd.Altitude).
				AddTag("latitude", sd.Latitude).
				AddTag("longitude", sd.Longitude).
				AddTag("unit", "m/s").
				AddField("avg15", wsa).
				AddField("max", wsm).
				SetTime(t)
			points = append(points, p)
		}
	}

	// Humidity
	if h, err := strconv.ParseFloat(sd.AirRelHumidityAvg, 32); err == nil && !math.IsNaN(h) {
		p := influxdb2.NewPointWithMeasurement(mHumidity).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "percent").
			AddField("avg15", h).
			SetTime(t)
		points = append(points, p)
	}

	// Precipitations.
	if pa, err := strconv.ParseFloat(sd.PrecipRtNrtTot, 32); err == nil && !math.IsNaN(pa) {
		p := influxdb2.NewPointWithMeasurement(mPrecipitations).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "mm").
			AddField("avg15", pa).
			SetTime(t)
		points = append(points, p)
	}

	// Snow.
	if sh, err := strconv.ParseFloat(sd.SnowHeight, 32); err == nil && !math.IsNaN(sh) {
		p := influxdb2.NewPointWithMeasurement(mSnow).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "m").
			AddField("height", sh).
			SetTime(t)
		points = append(points, p)
	}

	err = writeAPI.WritePoint(context.Background(), points...)

	return err
}

/*
func (s *Store) Query() error {
	client := influxdb2.NewClient(s.url, s.token)
	queryAPI := client.QueryApi(s.org)

	result, err := queryAPI.Query(context.Background(), `from(bucket:"`+s.bucket+`")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "`+mTemperature+`")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}

	// Ensures background processes finishes
	client.Close()
}
*/