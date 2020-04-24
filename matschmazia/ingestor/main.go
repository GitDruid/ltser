// Ingestor is a REST service that waits for sensors data as flat JSON and store them in an InfluxDB database.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/db/influxdb2"
	"goex/ltser/matschmazia/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	url    string
	org    string
	bucket string
	token  string
	host   string
	port   string
)

var dataStore db.Writer

func init() {
	flag.StringVar(&url, "u", "", "Target url of InfluxDB instance.")
	flag.StringVar(&org, "o", "", "Target organization.")
	flag.StringVar(&bucket, "b", "", "Target bucket.")
	flag.StringVar(&token, "t", "", "Auth token.")
	flag.StringVar(&host, "h", "localhost", "Service ip.")
	flag.StringVar(&port, "p", "8000", "Service port.")
}

func main() {
	flag.Parse()

	if url == "" || org == "" || bucket == "" || token == "" || host == "" || port == "" {
		fmt.Fprintln(flag.CommandLine.Output(), "Missing or empty parameter.")
		flag.Usage()
		os.Exit(-1)
	}

	dataStore = influxdb2.NewStore(url, org, bucket, token)

	http.HandleFunc("/sensordata", sensorDataHandler)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func sensorDataHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: add check for POST only request.

	// Not safe implementation: just for testing purpose.
	// See: https://haisum.github.io/2017/09/11/golang-ioutil-readall/
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stdout, "An error occurred: %q.", err)
		fmt.Fprintf(w, "An error occurred: %q.", err) //TODO: Improve response.
		return
	}

	var reading models.RawData
	err = json.Unmarshal(body, &reading)
	if err != nil {
		fmt.Fprintf(os.Stdout, "An error occurred: %q.", err)
		fmt.Fprintf(w, "An error occurred: %q.", err) //TODO: Improve response.
		return
	}

	//fmt.Fprintf(os.Stdout, "Data arrived: %v\n", reading)
	fmt.Fprintf(os.Stderr, ".")

	//go dataStore.Save(reading) // This will saturate "InfluDB Cloud Free" limit.
	err = dataStore.Write(reading)

	if err != nil {
		log.Printf("An error occurred: %v", err)
	}
	// TODO: Manage response to the caller.
}
