package main

import (
	"flag"
	"fmt"
	ext "goex/ltser/extensions"
	"goex/ltser/matschmazia/db"
	"goex/ltser/matschmazia/db/influxdb2"
	"goex/ltser/matschmazia/models"
	"goex/ltser/stats"
	"os"
	"time"

	"github.com/gitdruid/adf"
	//"github.com/berkmancenter/adf"
)

var (
	url     string
	org     string
	bucket  string
	token   string
	from    ext.DateTimeFlag
	to      ext.DateTimeFlag
	station string
)

var dataStore db.Reader

func init() {
	flag.StringVar(&url, "u", "", "Target url of InfluxDB instance.")
	flag.StringVar(&org, "o", "", "Target organization.")
	flag.StringVar(&bucket, "b", "", "Target bucket.")
	flag.StringVar(&token, "t", "", "Auth token.")
	flag.Var(&from, "from", "Start time (RFC3339 format).")
	flag.Var(&to, "to", "Finish time (RFC3339 format).")
	flag.StringVar(&station, "station", "", "Sensor station.")
}

func main() {
	flag.Parse()

	if url == "" || org == "" || bucket == "" || token == "" || from.Value().IsZero() || station == "" {
		fmt.Fprintln(flag.CommandLine.Output(), "Missing or empty parameter.")
		flag.Usage()
		os.Exit(-1)
	}

	if to.Value().IsZero() {
		to.Set(time.Now().Format(time.RFC3339))
	}

	dataStore = influxdb2.NewStore(url, org, bucket, token)

	res, err := dataStore.ReadAll(models.Snow, from.Value(), to.Value(), station)
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

	fixed, idx, err := stats.Hampel(res, 10, 5)

	for _, i := range idx {
		fmt.Printf("Value #%v (originally %g) was replaced by %g.\n", i, res[i], fixed[i])
	}
	fmt.Printf("Total fixed: %v\n", len(idx))

	// // TEST DEI CAMPIONI FIXATI
	// test, err = adf.New(fixed, 0, -1)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "An error occurred: %q.\n", err)
	// 	os.Exit(2)
	// }

	// test.Run()
	// fmt.Printf("Values in the series: %v\n", len(res))
	// fmt.Printf("Is stationary: %v\n", test.IsStationary())
}