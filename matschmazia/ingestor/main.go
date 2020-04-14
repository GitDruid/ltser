// Consumer is a REST service that waits for sensors data as flat JSON and store them in an InfluxDB database.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// More information on: https://browser.lter.eurac.edu/p/info.md
const (
	time              = iota // Date/time of measurement (UTC +1).
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

	body, err := ioutil.ReadAll(r.Body) // Not safe implementation: just for testing purpose.
	if err == nil {
		fmt.Fprintf(w, "Body = %q\n", body)
		log.Printf("Body = %s", body)
	}
}
