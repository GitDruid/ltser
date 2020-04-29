// Package influxdb2 provide an implementation of the db.ReadWriter interface for
// InfluxDB v2.0 databases.
package influxdb2 // import "goex/ltser/matschmazia/db/influxdb2"

import (
	"context"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/models"
	"math"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// A Store save and read data to and from an InfluxDB database.
type Store struct {
	url    string
	org    string
	bucket string
	token  string
}

// NewStore returns a new InfluxDB Store.
func NewStore(url, org, bucket, token string) *Store {
	influxDbStore := new(Store)
	influxDbStore.url = url
	influxDbStore.org = org
	influxDbStore.bucket = bucket
	influxDbStore.token = token

	return influxDbStore
}

// Save parse raw sensors' data and store valid data into separate measurements.
func (s *Store) Write(sd models.RawData) error {

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
		p := influxdb2.NewPointWithMeasurement(models.Temperature.Name()).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "celsius").
			AddField("avg15", temp).
			SetTime(t)
		points = append(points, p)
	}

	// Wind Speed.
	wsa, err := strconv.ParseFloat(sd.WindSpeedAvg, 32)
	if err != nil || math.IsNaN(wsa) {
		wsa, err = strconv.ParseFloat(sd.WindSpeed, 32) // In same cases values are in sd.WindSpeed, in others in sd.WindSpeedAvg.
	}
	if err == nil && !math.IsNaN(wsa) {
		p := influxdb2.NewPointWithMeasurement(models.WindSpeed.Name()).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "m/s").
			AddField("avg15", wsa).
			SetTime(t)
		points = append(points, p)

	}

	// Wind Gust.
	if wsm, err := strconv.ParseFloat(sd.WindSpeedMax, 32); err == nil && !math.IsNaN(wsm) {
		p := influxdb2.NewPointWithMeasurement(models.WindGust.Name()).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "m/s").
			AddField("max", wsm).
			SetTime(t)
		points = append(points, p)
	}

	// Humidity
	if h, err := strconv.ParseFloat(sd.AirRelHumidityAvg, 32); err == nil && !math.IsNaN(h) {
		p := influxdb2.NewPointWithMeasurement(models.Humidity.Name()).
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
		p := influxdb2.NewPointWithMeasurement(models.Precipitations.Name()).
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
		p := influxdb2.NewPointWithMeasurement(models.Snow.Name()).
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
// ReadAll data from a given measurement, time interval and station.
func (s *Store) ReadAll(m models.Measurement, rStart, rStop, station string) (res []float64, err error) {
	res = make([]float64, 0)
	client := influxdb2.NewClient(s.url, s.token)
	defer client.Close() // Ensures background processes finishes.

	queryAPI := client.QueryApi(s.org)

	result, err := queryAPI.Query(context.Background(),
		fmt.Sprintf(
			`from(bucket:%q)
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r._measurement == %q and r.station == %q)`, s.bucket, rStart, rStop, m, station))

	if err == nil {
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			r := result.Record()
			res = append(res, r.Value().(float64))
			fmt.Printf("row: %s\n", r.String())
		}
		if result.Err() != nil {
			err = result.Err()
		}
	}

	return
}
*/

// ReadAll data from a given measurement, time interval and station.
// In case of errors, the slice will contains data eventually read before the error happens.
func (s *Store) ReadAll(m models.Measurement, rStart, rStop, station string) ([]float64, error) {

	result, err := s.Read(m, rStart, rStop, station)
	if err != nil {
		return nil, err
	}

	series := make([]float64, 0)

	for {
		val, e := result.Next()
		if e == ErrEndOfRecords {
			break
		}

		if e != nil {
			err = e
			break
		}

		series = append(series, val)
	}

	return series, err
}

// Read returns a Result that needs to be iterated to obtain
// data from a given measurement, time interval and station.
func (s *Store) Read(m models.Measurement, rStart, rStop, station string) (db.Iterator, error) {
	client := influxdb2.NewClient(s.url, s.token)
	defer client.Close() // Ensures background processes finishes.

	queryAPI := client.QueryApi(s.org)

	qr, err := queryAPI.Query(context.Background(),
		fmt.Sprintf(
			`from(bucket:%q)
			|> range(start: %s, stop: %s) 
			|> filter(fn: (r) => r._measurement == %q and r.station == %q)`, s.bucket, rStart, rStop, m, station))

	if err != nil {
		return nil, err
	}

	var r = Result{queryResult: qr}
	return db.Iterator(r), err
}
