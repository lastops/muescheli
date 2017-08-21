package main

import (
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"bytes"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dutchcoders/go-clamd"
)

type FileResult struct {
	Filename	string
	Result		string
}

type ScanResult []FileResult

type App struct {
	Clam	*clamd.Clamd
	Router	*mux.Router
}

func (a *App) Initialize(clamdAddress string) {
	a.Clam = clamd.NewClamd(clamdAddress)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	var handler http.Handler

	handler = handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{http.MethodGet, http.MethodPost, http.MethodPut}))(a.Router)

	handler = handlers.ProxyHeaders(handler)
	handler = handlers.CompressHandler(handler)

	os.Stdout.WriteString("muescheli is ready and available on port " + strings.TrimPrefix(addr, ":"))
	log.Fatal(http.ListenAndServe(addr, handler))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/scan", auth(a.scanHttpUrl)).Methods(http.MethodGet)
	a.Router.HandleFunc("/scan", auth(a.scanBody)).Methods(http.MethodPut)
	a.Router.HandleFunc("/scan", auth(a.scanMultipart)).Methods(http.MethodPost)
}

func checkCredentials(username string, password string) bool {
	return true
}

func auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if !checkCredentials(user, pass) {
			http.Error(w, "Unauthorized.", 401)
			return
		}

		handler(w, r)
	}
}

func (a *App) scanMultipart(w http.ResponseWriter, r *http.Request) {
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
				result := a.scan(part)
				// write result
				fileResult := FileResult{part.FileName(), result }
				scanResult = append(scanResult, fileResult)
				fmt.Printf("scanned: %v, %v\n", part.FileName(), result)
			}
		}
	}

	respondWithJSON(w, http.StatusOK, scanResult)
}

func (a *App) scanBody(w http.ResponseWriter, r *http.Request) {
	scanResult := ScanResult{}

	buf, _ := ioutil.ReadAll(r.Body)
	part := ioutil.NopCloser(bytes.NewBuffer(buf))
	result := a.scan(part)
	// write result
	fileResult := FileResult{ "request body", result }
	scanResult = append(scanResult, fileResult)
	fmt.Printf("scanned: %v, %v\n", "request body", result)

	respondWithJSON(w, http.StatusOK, scanResult)
}

func (a *App) scanHttpUrl(w http.ResponseWriter, r *http.Request) {
	respondWithJSONError(w, http.StatusNotImplemented,"not implemented", "feature still to come...")
}

func (a *App) scan(r io.Reader) string {
	var abort chan bool
	response, _ := a.Clam.ScanStream(r, abort)

	for s := range response {
		return s.Status
	}
	return "ERROR"
}

func respondWithUserError(w http.ResponseWriter, description string) {
	respondWithJSONError(w, http.StatusBadRequest, "", description)
}

func respondWithServerError(w http.ResponseWriter, error error) {
	respondWithJSONError(w, http.StatusInternalServerError, "", error.Error())
}

func respondWithJSONError(w http.ResponseWriter, code int, error string, description string) {
	payload := map[string]string{}

	if len(error) > 0 {
		payload["error"] = error
	}

	if len(description) > 0 {
		payload["description"] = description
	}

	if len(payload) == 0 {
		payload = nil
	}

	respondWithJSON(w, code, payload)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload == nil {
		w.Write([]byte("{}"))
	} else {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}