// Pusher reads data from a csv file, transform each rows in a flat JSON
// and send them to StdOut or post them to a REST service.
//
// TODO: use goroutines and a buffered channel to decouple (and make concurrent) the reading logic from the sending logic.
package main

import (
	"encoding/csv"
	"flag"
	"goex/ltser/csvjson"
	"goex/ltser/sender"
	httpsender "goex/ltser/sender/http"
	stdoutsender "goex/ltser/sender/stdout"
	"io"
	"log"
	"os"
)

// Default values for parameters.
const (
	defFilename    = "./data.csv"
	defHeadersRows = 1
	noRowsLimit    = -1
	noURL          = ""
)

var (
	filename    string
	headersRows uint
	rowsToRead  int
	targetURL   string
	dataSender  sender.Sender
)

func init() {
	flag.StringVar(&filename, "f", defFilename, "Data .CSV file name.")
	flag.UintVar(&headersRows, "h", defHeadersRows, "Number of headers rows. First one is taken, the others are skipped.")
	flag.IntVar(&rowsToRead, "m", noRowsLimit, "Number of rows to read. Use -1 for no rows limit.")
	flag.StringVar(&targetURL, "u", noURL, "Target URL. If empty string, data are logged on StdOut.")
}

func main() {
	flag.Parse()

	// Open .CSV file.
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	csvRdr := csv.NewReader(f)
	csvRdr.ReuseRecord = true

	jsonRdr := csvjson.NewReader(*csvRdr)
	jsonRdr.HeadersRows = headersRows

	if targetURL == noURL {
		dataSender = stdoutsender.NewSender()
		jsonRdr.IndentFormat = true
		jsonRdr.Intent = "   "
	} else {
		dataSender = httpsender.NewSender(targetURL)
		jsonRdr.IndentFormat = false
	}

	for i := 0; rowsToRead < 0 || i < rowsToRead; i++ {
		// Read data.
		jsonBytes, err := jsonRdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		// Send data.
		//go dataSender.Send(jsonBytes) // This will saturate "InfluDB Cloud Free" limit.
		err = dataSender.Send(jsonBytes)
		if err != nil {
			log.Fatalf("An error occurred on row %v: %q. Aborting.", i, err)
		}
	}

	log.Print("Finished!")
}
