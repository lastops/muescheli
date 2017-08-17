package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/dutchcoders/go-clamd"
)

type ScanResult []FileResult

type FileResult struct {
	Filename	string
	Result		string
}

var clamAddress string

func handler(w http.ResponseWriter, r *http.Request) {
	// only do something with posts
	switch r.Method {
	case "POST":
		scanResult := ScanResult{}

		// if multipart form than get files from form
		reader, err := r.MultipartReader()
		if err == nil {
			// copy parts
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				// check if file and scan
				if part.FileName() != "" {
					result := scan(part)
					// write result
					fileResult := FileResult{ part.FileName(), result }
					scanResult = append(scanResult, fileResult)
					fmt.Printf("scanned: %v, %v\n", part.FileName(), result)
				}
			}
		} else { // just scan the whole body if not multipart form
			buf, _ := ioutil.ReadAll(r.Body)
			part := ioutil.NopCloser(bytes.NewBuffer(buf))
			result := scan(part)
			// write result
			fileResult := FileResult{ "request body", result }
			scanResult = append(scanResult, fileResult)
			fmt.Printf("scanned: %v, %v\n", "request body", result)
		}
		w.Header().Set("Content-Type", "application/json")
		// create json
		json, err := json.Marshal(scanResult)
		w.Write(json)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func scan(r io.Reader) string {
	c := clamd.NewClamd(clamAddress)

	var abort chan bool
	response, _ := c.ScanStream(r, abort)

	for s := range response {
		return s.Status
	}
	return "ERROR"
}

func main() {
	// muescheli port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8091"
	}

	// get clamd address
	clamd_port := os.Getenv("CLAMD_PORT")
	if len(clamd_port) == 0 {
		clamd_port = "3310"
	}
	clamd_host := os.Getenv("CLAMD_HOST")
	if len(clamd_host) == 0 {
		clamd_host = "localhost"
	}
	clamAddress = "tcp://" + clamd_host + ":" + clamd_port

	// wait for clamd
	test := clamd.NewClamd(clamAddress)
	fmt.Printf("waiting for clamd on %v\n", clamAddress)
	// TODO remove this hack and make nice function to handle resiliency not only on startup
	// use Version() and not Ping() because of nil pointer dereference that Ping() can cause
	for _, err := test.Version(); err != nil; _, err = test.Version() {
		fmt.Printf(".")
		time.Sleep(time.Second)
	}
	fmt.Printf("\nconnection to clamd successful on %v\n", clamAddress)

	// define handler
	http.HandleFunc("/scan", handler)
	// listen on port
	http.ListenAndServe(":" + port, nil)
}