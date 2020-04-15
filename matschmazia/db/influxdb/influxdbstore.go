// Package influxdb provide an implementation of the store interface to save data to database.
package influxdb // import "goex/ltser/matschmazia/db/influxdb"

import (
	"context"
	"fmt"
	"goex/ltser/matschmazia/models"
	"math"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// A Store save data to database.
type Store struct {
}

// NewStore returns a new Store.
func NewStore() *Store {
	influxDbStore := new(Store)

	return influxDbStore
}

// Save store data in InfluxDB measurements.
func (s *Store) Save(sd models.SensorData) error {
	orgName := "galassiasoft.com"
	bucketName := "ltser-bucket"

	mTemperature := "temperature"
	mWind := "wind"
	mHumidity := "humidity"
	mPrecipitations := "precipitations"

	client := influxdb2.NewClient("https://eu-central-1-1.aws.cloud2.influxdata.com", "NCWF7CXKdcoOGJ-dOA4EEIl-OaTZGUZLw6cEtlTho8nI7J-iznobcrs94W8jMZjLjyPN9NX8O48iPGN6-Aq18Q==")

	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteApi(orgName, bucketName)

	// Temperature.
	if t, err := strconv.ParseFloat(sd.AirTempAvg, 32); err == nil && !math.IsNaN(t) {
		point := influxdb2.NewPointWithMeasurement(mTemperature).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "celsius").
			AddField("avg15", t).
			SetTime(time.Now())
		writeAPI.WritePoint(point)
	}

	// Wind.
	wsa, err := strconv.ParseFloat(sd.WindSpeedAvg, 32)
	if err != nil || math.IsNaN(wsa) {
		wsa, err = strconv.ParseFloat(sd.WindSpeed, 32) // In same cases values are in sd.WindSpeed, in others in sd.WindSpeedAvg.
	}
	if err == nil && !math.IsNaN(wsa) {
		if wsm, err := strconv.ParseFloat(sd.WindSpeedMax, 32); err == nil && !math.IsNaN(wsm) {
			point := influxdb2.NewPointWithMeasurement(mWind).
				AddTag("station", sd.Station).
				AddTag("altitude", sd.Altitude).
				AddTag("latitude", sd.Latitude).
				AddTag("longitude", sd.Longitude).
				AddTag("unit", "m/s").
				AddField("avg15", wsa).
				AddField("max", wsm).
				SetTime(time.Now())
			writeAPI.WritePoint(point)
		}
	}

	// Humidity
	if h, err := strconv.ParseFloat(sd.AirRelHumidityAvg, 32); err == nil && !math.IsNaN(h) {
		point := influxdb2.NewPointWithMeasurement(mHumidity).
			AddTag("station", sd.Station).
			AddTag("altitude", sd.Altitude).
			AddTag("latitude", sd.Latitude).
			AddTag("longitude", sd.Longitude).
			AddTag("unit", "percent").
			AddField("avg15", h).
			SetTime(time.Now())
		writeAPI.WritePoint(point)
	}

	// Precipitations.
	if p, err := strconv.ParseFloat(sd.PrecipRtNrtTot, 32); err == nil && !math.IsNaN(p) {
		if sh, err := strconv.ParseFloat(sd.SnowHeight, 32); err == nil && !math.IsNaN(sh) {
			point := influxdb2.NewPointWithMeasurement(mPrecipitations).
				AddTag("station", sd.Station).
				AddTag("altitude", sd.Altitude).
				AddTag("latitude", sd.Latitude).
				AddTag("longitude", sd.Longitude).
				AddTag("unit", "mm").
				AddField("avg15", p).
				AddField("show_height", sh*1000). // Scaled from m to mm.
				SetTime(time.Now())
			writeAPI.WritePoint(point)
		}
	}

	writeAPI.Flush()

	// Ensures background processes finishes
	client.Close()

	return nil
}

func query() {
	orgName := "galassiasoft.com"
	bucketName := "ltser-bucket"
	measurementName := "test-measurement"

	// create new client with default option for server url authenticate by token
	client := influxdb2.NewClient("https://eu-central-1-1.aws.cloud2.influxdata.com", "NCWF7CXKdcoOGJ-dOA4EEIl-OaTZGUZLw6cEtlTho8nI7J-iznobcrs94W8jMZjLjyPN9NX8O48iPGN6-Aq18Q==")

	// get query client
	queryAPI := client.QueryApi(orgName)

	// get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"`+bucketName+`")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "`+measurementName+`")`)
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
