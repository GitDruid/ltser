// Pusher reads data from a csv file, transform each row in a flat JSON
// and send them to StdOut or post them to a REST service.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"goex/ltser/csvjson"
	ext "goex/ltser/extensions"
	"goex/ltser/sender"
	httpsender "goex/ltser/sender/http"
	stdoutsender "goex/ltser/sender/stdout"
	"io"
	"log"
	"os"
	"time"
)

// Default values for parameters.
const (
	defFilename       = "./data.csv"
	defHeadersRows    = 1
	noRowsLimit       = -1
	noURL             = ""
	defBufferSize     = 1
	defMaxConcurrency = 1
)

type task byte

const (
	readerTask task = 0
	senderTask task = 1
)

type controlMsg struct {
	err     error
	isFatal bool
	origin  task
}

var (
	filename       string
	headersRows    uint
	rowsToRead     int
	targetURL      string
	bufferSize     ext.NotZeroUint32
	maxConcurrency ext.NotZeroUint32
	jsonRdr        csvjson.Reader
	dataSender     sender.Sender
	chData         chan []byte
	chControl      chan controlMsg
	onGoing        int
	msg            controlMsg
)

func init() {
	flag.StringVar(&filename, "f", defFilename, "Data .CSV file name.")
	flag.UintVar(&headersRows, "h", defHeadersRows, "Number of headers rows. First one is taken, the others are skipped.")
	flag.IntVar(&rowsToRead, "m", noRowsLimit, "Number of rows to read. Use -1 for no rows limit.")
	flag.StringVar(&targetURL, "u", noURL, "Target URL. If empty string, data are logged on StdOut.")
	flag.Var(&bufferSize, "b", "Buffer size while reading.")
	flag.Var(&maxConcurrency, "c", "Max concurrency. If greater than 1, sequential data processing is not guaranteed.")
}

func main() {
	defer trace("pusher")()

	flag.Parse()

	// Open .CSV file.
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v", err)
		os.Exit(1)
	}

	csvRdr := csv.NewReader(f)
	csvRdr.ReuseRecord = true

	jsonRdr = *csvjson.NewReader(*csvRdr)
	jsonRdr.HeadersRows = headersRows

	if targetURL == noURL {
		dataSender = stdoutsender.NewSender()
		jsonRdr.IndentFormat = true
		jsonRdr.Indent = "   "
	} else {
		dataSender = httpsender.NewSender(targetURL)
		jsonRdr.IndentFormat = false
	}

	chData = make(chan []byte, bufferSize.Value())
	chControl = make(chan controlMsg)

	go read()
	for i := uint32(0); i < maxConcurrency.Value(); i++ {
		go send()
	}

	for {
		msg = <-chControl

		if msg.origin == readerTask && msg.err == io.EOF {
			break
		}
		if msg.origin == readerTask {
			onGoing++
		}
		if msg.origin == senderTask {
			onGoing--
		}
		logMsg(msg)
	}

	// Reader finished. Waiting for senders to complete any pending tasks.
	for onGoing > 0 {
		msg = <-chControl
		logMsg(msg)
		onGoing--
	}

	fmt.Fprintln(os.Stderr, "\nFinished!")
}

func read() {
	for i := 0; rowsToRead < 0 || i < rowsToRead; i++ {
		jsonBytes, err := jsonRdr.Read()
		if err == io.EOF {
			break
		}
		if err == nil {
			chData <- jsonBytes
		}
		chControl <- controlMsg{err: err, origin: readerTask}
	}

	chControl <- controlMsg{err: io.EOF, origin: readerTask}
}

func send() {
	for i := 0; ; i++ {
		jsonBytes := <-chData
		err := dataSender.Send(jsonBytes)
		if err != nil {
			chControl <- controlMsg{err: err, isFatal: true, origin: senderTask}
		} else {
			chControl <- controlMsg{origin: senderTask}
		}
	}
}

func trace(message string) func() {
	start := time.Now()
	log.Printf("enter %s", message)
	return func() { log.Printf("exit %s (%s)", message, time.Since(start)) }
}

func logMsg(msg controlMsg) {
	switch {
	case msg.err == nil && msg.origin == readerTask:
		fmt.Fprintf(os.Stderr, "r")
	case msg.err == nil && msg.origin == senderTask:
		fmt.Fprintf(os.Stderr, "s")
	case msg.isFatal:
		fmt.Fprintf(os.Stderr, "\nAn error occurred (%s). Aborted.", msg.err)
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "\nAn error occurred (%s).", msg.err)
	}
}
