// Pusher reads data from a csv file, transform each rows in a flat JSON
// and send them to StdOut or post them to a REST service.
//
// TODO 2: add a flag to switch between sequential implementation (that guarantees
// rows order) and parallel implementation (no order guaranteed).
package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"goex/ltser/csvjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	send        func(int, []byte) // Function variable to change behavior based on targetURL.
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
		send = sendToStdOut
		jsonRdr.IndentFormat = true
		jsonRdr.Intent = "   "
	} else {
		send = sendToTargetURL
		jsonRdr.IndentFormat = false
	}

	for i := 0; rowsToRead < 0 || i < rowsToRead; i++ {
		// Read data.
		jsonBytes, err := jsonRdr.Read()
		if err != nil {
			if err == io.EOF {
				log.Print("Finished!")
				os.Exit(0)
			}
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		// Send data.
		send(i, jsonBytes)
	}

	log.Print("Finished!")
}

func sendToStdOut(i int, b []byte) {
	fmt.Printf("%s\n", b)
}

func sendToTargetURL(i int, b []byte) {
	r, err := http.Post(targetURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatalf("An error occurred on row %v: %v. Aborting.", i, err)
	}
	_, err = ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Fatalf("An error occurred on row %v: %v. Aborting.", i, err)
	}
	if r.StatusCode == http.StatusOK {
		log.Print(".")
	} else {
		log.Printf("An error occurred on row %v: response status %q.", i, r.Status)
	}
}
