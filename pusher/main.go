// Pusher reads data from a csv file, transform each row in a flat JSON
// and send them to StdOut or post them to a REST service.
package main

import (
	"encoding/csv"
	"errors"
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
	line    uint
}

type dataMsg struct {
	data []byte
	line uint
}

var (
	filename        string
	headersRows     uint
	rowsToRead      int
	targetURL       string
	bufferSize      ext.NotZeroUint32
	maxConcurrency  ext.NotZeroUint32
	jsonRdr         csvjson.Reader
	dataSender      sender.Sender
	chData          chan dataMsg
	chControl       chan controlMsg
	errEndOfSending = errors.New("End Of Sending")
	endOfSendingMsg = controlMsg{err: errEndOfSending, isFatal: false, origin: senderTask, line: 0}
	totalLines      uint
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

	chData = make(chan dataMsg, bufferSize.Value())
	chControl = make(chan controlMsg)

	go read()
	for i := uint32(0); i < maxConcurrency.Value(); i++ {
		go send()
	}

	for {
		msg := <-chControl

		if msg.origin == readerTask && msg.err == io.EOF {
			totalLines = msg.line
			break
		}

		logMsg(msg)
	}

	close(chData)

	// Reader finished. Waiting for senders to complete any pending tasks.
	for i := uint32(0); i < maxConcurrency.Value(); {
		msg := <-chControl
		if msg == endOfSendingMsg {
			i++
		} else {
			logMsg(msg)
		}
	}

	//TODO: sender fatal error. Stop reader and wait for results from other senders.

	close(chControl)

	fmt.Fprintf(os.Stderr, "\nFinished processing %v lines.\n", totalLines)
}

func read() {
	i := uint(0)
	for i = 0; rowsToRead < 0 || i < uint(rowsToRead); i++ { // Cast only if >= 0.
		jsonBytes, err := jsonRdr.Read()
		if err == io.EOF {
			break
		}
		if err == nil {
			chData <- dataMsg{data: jsonBytes, line: i}
		}
		chControl <- controlMsg{err: err, origin: readerTask, line: i}
	}

	chControl <- controlMsg{err: io.EOF, origin: readerTask, line: i}
}

func send() {
	for {
		msg, more := <-chData
		if !more {
			chControl <- endOfSendingMsg
			return
		}

		err := dataSender.Send(msg.data)
		fatal := false
		if err != nil {
			fatal = true // TODO: Add error analysis logic here.
		}
		chControl <- controlMsg{err: err, isFatal: fatal, origin: senderTask, line: msg.line}
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
		fmt.Fprintf(os.Stderr, "\nAn error occurred on line %v (%s). Aborted.", msg.line, msg.err)
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "\nAn error occurred on line %v (%s).", msg.line, msg.err)
	}
}
