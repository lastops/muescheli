package main

import (
	"io"
	"os"
	"fmt"
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
	// get clamd address
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3310"
	}
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	clamAddress = "tcp://" + host + ":" + port

	// test if clamd is around
	test := clamd.NewClamd(clamAddress)
	ping := test.Ping()
	if ping != nil {
		fmt.Printf("can't connect to clamd on %v\n", clamAddress)
		os.Exit(1)
	}
	fmt.Printf("connection to clamd successful on %v\n", clamAddress)

	// define handler
	http.HandleFunc("/scan", handler)
	// listen on port
	http.ListenAndServe(":8091", nil)
}