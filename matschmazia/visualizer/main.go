package main

import (
	"flag"
	"fmt"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/db/influxdb2"
	"goex/ltser/matschmazia/models"
	"os"

	"github.com/gitdruid/adf"
	//"github.com/berkmancenter/adf"
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
	res, err := dataStore.Read(models.WindSpeed, "2020-03-20T00:00:00Z", "2020-04-20T23:59:00Z", "b1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %q.\n", err)
		os.Exit(1)
	}

	test, err := adf.New(res, 0, -1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %q.\n", err)
		os.Exit(2)
	}

	test.Run()

	fmt.Printf("Values in the series: %v\n", len(res))
	fmt.Printf("Is stationary: %v\n", test.IsStationary())
}
