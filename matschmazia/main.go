// Eurac reads sensors' data from a csv file
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	defFilename   = "./data.csv"
	defHeaderRows = 2
	defRowsToRead = -1
)

const (
	time = iota
	station
	landuse
	altitude
	latitude
	longitude
	airRelHumidityAvg
	airTempAvg
	nrUpSwAvg
	precipRtNrtTot
	snowHeight
	swRadiationAvg
	windDir
	windSpeed
	windSpeedAvg
	windSpeedMax
)

var filename string
var headerRows int
var rowsToRead int

func init() {
	flag.StringVar(&filename, "f", defFilename, "Sensor data .CSV file name.")
	flag.IntVar(&headerRows, "h", defHeaderRows, "Number of header rows to skip.")
	flag.IntVar(&rowsToRead, "m", defRowsToRead, "Number of rows to read.")
}

func main() {

	flag.Parse()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	r := csv.NewReader(f)
	r.ReuseRecord = true

	for i := 0; rowsToRead < 0 || i < rowsToRead+headerRows; i++ {
		record, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("An error occurred: %v", err)
		}

		if i < headerRows {
			continue
		}

		fmt.Println(record)
	}

}
