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
	_ "net/http/pprof"
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

	// Profiling service.
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	dataStore = influxdb2.NewStore(url, org, bucket, token)

	http.HandleFunc("/sensordata", sensorDataHandler)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func sensorDataHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: add check for POST only request.

	// Not safe implementation: just for testing purpose.
	// See: https://haisum.github.io/2017/09/11/golang-ioutil-readall/
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("An error occurred: %q.\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var reading models.RawData
	err = json.Unmarshal(body, &reading)
	if err != nil {
		log.Printf("An error occurred: %q.\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(".")

	err = dataStore.Write(reading)
	if err != nil {
		log.Printf("An error occurred: %q.", err)
		http.Error(w, "unable to save data", http.StatusInternalServerError)
		return
	}

	http.Error(w, "success", http.StatusOK)
}
