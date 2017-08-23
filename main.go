package main

import (
	"os"
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
	// run app
	a.Run(":" + port)
}