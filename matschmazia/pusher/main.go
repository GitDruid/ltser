// Pusher reads data from a csv file, transform each row in a flat JSON
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
)

func init() {
	flag.StringVar(&filename, "f", defFilename, "Data .CSV file name.")
	flag.UintVar(&headersRows, "h", defHeadersRows, "Number of headers rows. First one is taken, the others are skipped.")
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
	for i := uint(0); i < headersRows; i++ {
		record := readFrom(r)
		if i == 0 {
			headers = append(headers, record...) // Cloning values since record is a slice and r.ReuseRecord = true.
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

		jsonBytes, err := json.MarshalIndent(m, "", "   ")
		if err != nil {
			log.Printf("Skipped malformed row #%v (%s).", i, err)
			continue
		}

		//TODO: refactor in a separate parametric function.
		if targetURL == noURL {
			fmt.Printf("%s\n", jsonBytes)
		} else {
			r, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBytes))
			if err != nil {
				log.Fatalf("An error occurred on row %v: %v. Aborting.", i, err)
			}
			b, err := ioutil.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				log.Fatalf("An error occurred on row %v: %v. Aborting.", i, err)
			}
			if r.StatusCode == http.StatusOK {
				log.Printf("%s\n", b)
			} else {
				log.Printf("An error occurred on row %v: response status %q.", i, r.Status)
			}
		}
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
		return nil, errors.New("keys and values sizes dont match")
	}

	m := make(map[string]string)

	for i := 0; i < len(k); i++ {
		m[k[i]] = v[i]
	}

	return m, nil
}
