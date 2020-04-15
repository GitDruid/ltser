// Ingestor is a REST service that waits for sensors data as flat JSON and store them in an InfluxDB database.
package main

import (
	"encoding/json"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/db/influxdb"
	"goex/ltser/matschmazia/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var dataStore db.Store

func main() {
	dataStore = influxdb.NewStore()

	http.HandleFunc("/sensordata", sensorDataHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
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

	var reading models.SensorData
	err = json.Unmarshal(body, &reading)
	if err != nil {
		fmt.Fprintf(os.Stdout, "An error occurred: %q.", err)
		fmt.Fprintf(w, "An error occurred: %q.", err) //TODO: Improve response.
		return
	}

	fmt.Fprintf(os.Stdout, "Data arrived: %v", reading)

	dataStore.Save(reading)

	// TODO: Manage response to the caller.
}
