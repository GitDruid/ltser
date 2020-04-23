package main

import (
	"flag"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/db/influxdb2"
	"goex/ltser/matschmazia/models"
	"os"
)

var (
	url    string
	org    string
	bucket string
	token  string
)

var dataStore db.Store

func init() {
	flag.StringVar(&url, "u", "", "Target url of InfluxDB instance.")
	flag.StringVar(&org, "o", "", "Target organization.")
	flag.StringVar(&bucket, "b", "", "Target bucket.")
	flag.StringVar(&token, "t", "", "Auth token.")
}

func main() {
	flag.Parse()

	if url == "" || org == "" || bucket == "" || token == "" {
		fmt.Fprintln(flag.CommandLine.Output(), "Missing or empty parameter.")
		flag.Usage()
		os.Exit(-1)
	}

	dataStore = influxdb2.NewStore(url, org, bucket, token)

	//err := dataStore.Read(models.Temperature, "-15h", "now()")
	err := dataStore.Read(models.Temperature, "2020-04-01T00:00:00Z", "2020-04-02T00:00:00Z")
	if err != nil {
		fmt.Printf("An error occurred: %q", err)
	}
}
