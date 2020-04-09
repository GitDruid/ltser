// Pusher reads data from a csv file, transform each rows in a flat JSON
// and send them to StdOut or post them to a REST service.
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	send        func(i int, b []byte)
)

func init() {
	flag.StringVar(&filename, "f", defFilename, "Data .CSV file name.")
	flag.UintVar(&headersRows, "h", defHeadersRows, "Number of headers rows. First one is taken, the others are skipped.")
	flag.IntVar(&rowsToRead, "m", noRowsLimit, "Number of rows to read. Use -1 for no rows limit.")
	flag.StringVar(&targetURL, "u", noURL, "Target URL. If empty string, data are logged on StdOut.")
}

func main() {

	flag.Parse()

	// Function variable to change behavior based on targetURL.
	if targetURL == noURL {
		send = sendToStdOut
	} else {
		send = sendToTargetURL
	}

	// Open .CSV file.
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
	r := csv.NewReader(f)
	r.ReuseRecord = true

	// Read headers.
	headers := []string{}
	for i := uint(0); i < headersRows; i++ {
		record := readFrom(r)
		if i == 0 {
			headers = append(headers, record...) // Cloning values since record is a slice and r.ReuseRecord = true.
		}
	}

	for i := 0; rowsToRead < 0 || i < rowsToRead; i++ {
		// Read data.
		record := readFrom(r)

		// Trasform data.
		if i == 0 && len(headers) == 0 { // If headers are missing, generate default columns' names.
			for j := 0; j < len(record); j++ {
				headers = append(headers, "column"+strconv.Itoa(j))
			}
		}
		m, err := toMap(headers, record)
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}
		jsonBytes, err := json.MarshalIndent(m, "", "   ")
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		// Send data.
		send(i, jsonBytes)
	}

	log.Print("Finished!")
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
		return nil, errors.New("keys and values sizes don't match")
	}

	m := make(map[string]string)

	for i := 0; i < len(k); i++ {
		m[k[i]] = v[i]
	}

	return m, nil
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
