// Consumer is a REST service that waits for sensors data as flat JSON and store them in an InfluxDB database.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// More information on: https://browser.lter.eurac.edu/p/info.md
const (
	dateTime          = iota // [should be "time"] Date/time of measurement (UTC +1).
	station                  // Station code.
	landuse                  // me = meadows, pa = pasture, bs = bare soil, fo = forest
	altitude                 // Altitude of the station in meters.
	latitude                 // Latitude, coordinates in decimal degrees.
	longitude                // Longitude, coordinates in decimal degrees.
	airRelHumidityAvg        // Relative humidity in percent (15 min average).
	airTempAvg               // Air temperature in degree celsius (15 min average).
	nrUpSwAvg                // Undocumented.
	precipRtNrtTot           // Precipitation in mm (15 min cumulative sum).
	snowHeight               // Snow height in meter.
	srAvg                    // Global solar radiation in Watt square meter (15 min average).
	windDir                  // Wind direction in degrees (15 min average).
	windSpeed                // Undocumented.
	windSpeedAvg             // Wind speed in m/s.
	windSpeedMax             // Wind gust in m/s.
)

func main() {
	//writeReadTest()

	http.HandleFunc("/sensordata", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	//fmt.Scanln() // wait for Enter Key
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}

	body, err := ioutil.ReadAll(r.Body) // Not safe implementation: just for testing purpose. See: https://haisum.github.io/2017/09/11/golang-ioutil-readall/
	r.Body.Close()
	if err != nil {
		log.Print(err)
	} else {
		fmt.Fprintf(w, "Body = %q\n", body)
		log.Printf("Body = %s", body)
	}
}

func writeReadTest() {
	orgName := "galassiasoft.com"
	bucketName := "ltser-bucket"
	measurementName := "test-measurement"
	// create new client with default option for server url authenticate by token
	//client := influxdb2.NewClient("http://localhost:9999", "my-token")
	client := influxdb2.NewClient("https://eu-central-1-1.aws.cloud2.influxdata.com", "NCWF7CXKdcoOGJ-dOA4EEIl-OaTZGUZLw6cEtlTho8nI7J-iznobcrs94W8jMZjLjyPN9NX8O48iPGN6-Aq18Q==")

	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteApiBlocking(orgName, bucketName)

	// create point using full params constructor
	p := influxdb2.NewPoint(measurementName,
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)

	// create point using fluent style
	p = influxdb2.NewPointWithMeasurement(measurementName).
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)

	// Or write directly line protocol
	line := fmt.Sprintf("%s,unit=temperature avg=%f,max=%f", measurementName, 23.5, 45.0)
	writeAPI.WriteRecord(context.Background(), line)

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
