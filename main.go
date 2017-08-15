package main

import (
	"io"
	"os"
	"fmt"
	"strings"
	"net/http"
	"github.com/dutchcoders/go-clamd"
)

var clamAddress string = "tcp://localhost:3310"

func scan(w http.ResponseWriter, r *http.Request) {
	// only do something with posts
	if r.Method == "POST" {
		c := clamd.NewClamd(clamAddress)

		// multipart reader
		reader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		// copy parts
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			// check if file
			if part.FileName() != "" {

				var abort chan bool
				response, _ := c.ScanStream(part, abort)

				for s := range response {
					if strings.Contains(s.Status, "FOUND") {
						w.Write([]byte(part.FileName()))
					}

					fmt.Printf("scanned: %v, %v\n", part.FileName(), s.Status)
				}
			}
		}
	}
}

func main() {
	test := clamd.NewClamd(clamAddress)
	ping := test.Ping()
	if ping != nil {
		fmt.Printf("can't connect to clamd on %v\n", clamAddress)
		os.Exit(1)
	}
	fmt.Printf("connection to clamd successful on %v\n", clamAddress)

	// define handler
	http.HandleFunc("/scan", scan)
	// listen on port
	http.ListenAndServe(":8091", nil)
}