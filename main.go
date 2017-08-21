package main

import (
	"os"
	"fmt"
	"time"
)

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
	clamAddress := "tcp://" + clamd_host + ":" + clamd_port

	// init app
	a := App{}
	a.Initialize(clamAddress)

	// wait for clamd
	fmt.Printf("waiting for clamd on %v\n", clamAddress)
	// TODO remove this hack and make nice function to handle resiliency not only on startup
	// use Version() and not Ping() because of nil pointer dereference that Ping() can cause
	for _, err := a.Clam.Version(); err != nil; _, err = a.Clam.Version() {
		fmt.Printf(".")
		time.Sleep(time.Second)
	}
	fmt.Printf("\nconnection to clamd successful on %v\n", clamAddress)

	// run app
	a.Run(":" + port)
}