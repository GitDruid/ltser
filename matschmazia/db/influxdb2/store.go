// Package influxdb2 provide an implementation of the db.ReadWriter interface for
// InfluxDB v2.0 databases.
package influxdb2 // import "goex/ltser/matschmazia/db/influxdb2"

import (
	"context"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/models"
	"goex/ltser/timeseries"
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

// ReadAll data from a given measurement, time interval and station.
// In case of errors, the slice will contains data eventually read before the error happens.
func (s *Store) ReadAll(m models.Measurement, rStart, rStop time.Time, station string) (*models.Observations, error) {

	result, err := s.Read(m, rStart, rStop, station)
	if err != nil {
		return nil, err
	}

	o := models.Observations{
		Station:     result.Station(),
		Measurement: result.Measurement(),
		Measures:    timeseries.TimeSeries{}}

	for {
		val, e := result.Next()
		if e == ErrEndOfRecords {
			break
		}

		if e != nil {
			err = e
			break
		}

		o.Measures.AddWithTime(*val)
	}

	return &o, err
}

// Read returns a Result that needs to be iterated to obtain
// data from a given measurement, time interval and station.
func (s *Store) Read(m models.Measurement, rStart, rStop time.Time, station string) (db.ObservationsIterator, error) {
	client := influxdb2.NewClient(s.url, s.token)
	defer client.Close() // Ensures background processes finishes.

	queryAPI := client.QueryApi(s.org)

	qr, err := queryAPI.Query(context.Background(),
		fmt.Sprintf(
			`from(bucket:%q)
			|> range(start: %s, stop: %s) 
			|> filter(fn: (r) => r._measurement == %q and r.station == %q)`,
			s.bucket, rStart.Format(time.RFC3339), rStop.Format(time.RFC3339), m, station))

	if err != nil {
		return nil, err
	}

	var r = Result{queryResult: qr}

	// First record is read here. The others are read in the Iterator.
	if r.queryResult.Next() {
		i, _ := strconv.ParseInt(r.queryResult.Record().ValueByKey("altitude").(string), 10, 32)
		r.station.Altitude = int32(i)
		f, _ := strconv.ParseFloat(r.queryResult.Record().ValueByKey("latitude").(string), 32)
		r.station.Latitude = float32(f)
		f, _ = strconv.ParseFloat(r.queryResult.Record().ValueByKey("longitude").(string), 32)
		r.station.Longitude = float32(f)
		r.station.Name = r.queryResult.Record().ValueByKey("station").(string)
		r.measurement = m
		r.currentError = r.queryResult.Err()
		r.currentValue = &timeseries.TimeValue{
			Time:  r.queryResult.Record().Time(),
			Value: r.queryResult.Record().Value(),
		}
	} else {
		r.currentError = ErrEndOfRecords
		r.currentValue = nil
	}

	return db.ObservationsIterator(&r), err
}
