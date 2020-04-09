// Pusher reads data from a csv file and post them to a REST service in a flat JSON.
package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Default values for parameters.
const (
	defFilename   = "./data.csv"
	defHeaderRows = 1
	noRowsLimit   = -1
	noURL         = ""
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

var filename string
var headerRows uint
var rowsToRead int
var targetURL string

func init() {
	flag.StringVar(&filename, "f", defFilename, "Data .CSV file name.")
	flag.UintVar(&headerRows, "h", defHeaderRows, "Number of header rows. First one is taken, the others are skipped.")
	flag.IntVar(&rowsToRead, "m", noRowsLimit, "Number of rows to read. Use -1 for no rows limit.")
	flag.StringVar(&targetURL, "u", noURL, "Target URL. If empty string, data are logged on StdOut.")
}

func main() {

	flag.Parse()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	r := csv.NewReader(f)
	r.ReuseRecord = true

	// Read headers.
	headers := []string{}
	for i := uint(0); i < headerRows; i++ {
		record := readFrom(r)
		if i == 0 {
			headers = append(headers, record...) // Cloning values since record is a slice and r.ReuseRecord = true
			//headers = []string{"pippo", "pluto"}        // Test: few headers than data columns.
			//headers = append(headers, "pippo", "pluto") // Test: more headers than data columns.
		}
	}

	// Read data.
	for i := 0; rowsToRead < 0 || i < rowsToRead; i++ {
		record := readFrom(r)

		if i == 0 && len(headers) == 0 { // If headers are missing, generate default columns' names.
			for i := 0; i < len(record); i++ {
				headers = append(headers, "column"+strconv.Itoa(i))
			}
		}

		m, err := toMap(headers, record)
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		jsonStr := string(b)

		if targetURL == noURL {
			fmt.Println(jsonStr)
		} else {
			//TODO: POST to remote.
		}
	}
}

func readFrom(r *csv.Reader) (record []string) {
	record, err := r.Read()

	if err != nil {
		if err == io.EOF {
			log.Print("Finished!")
			os.Exit(0)
		}
		log.Fatalf("An error occurred: %v", err)
	}

	return record
}

func toMap(k []string, v []string) (map[string]string, error) {

	if len(v) != len(k) {
		return nil, errors.New("keys and values sizes dont match")
	}

	m := make(map[string]string)

	for i := 0; i < len(k); i++ {
		m[k[i]] = v[i]
	}

	return m, nil
}
