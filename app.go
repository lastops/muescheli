package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lastops/go-clamd"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FileResult struct {
	Filename string
	Result   string
}

type ScanResult []FileResult

type App struct {
	Clam   *clamd.Clamd
	Router *mux.Router
}

var graylogAddr = os.Getenv("GRAYLOGADDR")
var application = os.Getenv("APPLICATION")

func logger(source string, flavour string, result string) {

	if len(graylogAddr) == 0 {
		graylogAddr = "localhost:12201"
	}

	if len(application) == 0 {
		application = "av-dev"
	}

	// configure logging
	log.SetFormatter(&log.JSONFormatter{})

	// If using TCP
	gelfWriter, err := gelf.NewTCPWriter(graylogAddr)
	// If using UDP
	//gelfWriter, err := gelf.NewUDPWriter(graylogAddr)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}
	// log to both stderr and graylog2
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))

	log.WithFields(log.Fields{
		"application": application,
		"source":      source,
		"flavour":     flavour,
	}).Info(result)

}

func (a *App) Initialize(clamdAddress string) {

	a.Clam = clamd.NewClamd(clamdAddress)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	var handler http.Handler

	handler = a.Router
	handler = handlers.ProxyHeaders(handler)
	handler = handlers.CompressHandler(handler)

	handler = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}).Handler(handler)

	log.Printf("muescheli is ready and available on port %s", strings.TrimPrefix(addr, ":"))
	log.Fatal(http.ListenAndServe(addr, handler))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/scan", auth(a.scanMultipart)).Methods(http.MethodPost)
	a.Router.HandleFunc("/scan", auth(a.scanBody)).Methods(http.MethodPut)
	a.Router.Path("/scan").Queries("url", "{url}").HandlerFunc(auth(a.scanHttpUrl)).Methods(http.MethodGet)

	// endpoints to check if webservice is up
	a.Router.HandleFunc("/liveness", a.livenessCheck).Methods(http.MethodGet)
	a.Router.HandleFunc("/readiness", a.readinessCheck).Methods(http.MethodGet)
}

func checkCredentials(username string, password string) bool {
	// if env variables are empty or not set ignore credentials
	if user, isUserSet := os.LookupEnv("MUESCHELI_USER"); isUserSet && len(user) > 0 {
		if pass, isPassSet := os.LookupEnv("MUESCHELI_PASSWORD"); isPassSet && (pass != password || user != username) {
			return false
		}
	}
	return true
}

func auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if !checkCredentials(user, pass) {
			respondWithJSONError(w, http.StatusUnauthorized, "Unauthorized", "wrong credentials")
			return
		}

		handler(w, r)
	}
}

func (a *App) scanMultipart(w http.ResponseWriter, r *http.Request) {
	scanResult := ScanResult{}

	// read multipart
	reader, err := r.MultipartReader()
	defer r.Body.Close()
	if err != nil {
		respondWithUserError(w, err.Error())
		return
	}

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
			fileResult := FileResult{part.FileName(), result}
			scanResult = append(scanResult, fileResult)
			logger(part.FileName(), "file", result)
		}
	}

	respondWithJSON(w, http.StatusOK, scanResult)
}

func (a *App) scanBody(w http.ResponseWriter, r *http.Request) {
	scanResult := ScanResult{}

	buf, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithServerError(w, err)
		return
	}
	part := ioutil.NopCloser(bytes.NewBuffer(buf))
	result := a.scan(part)
	// write result
	fileResult := FileResult{"request body", result}
	scanResult = append(scanResult, fileResult)
	logger(string(buf), "request body", result)

	respondWithJSON(w, http.StatusOK, scanResult)
}

func (a *App) scanHttpUrl(w http.ResponseWriter, r *http.Request) {
	scanResult := ScanResult{}

	// get url parameter from request
	vars := mux.Vars(r)
	url := vars["url"]
	defer r.Body.Close()

	// download url
	response, err := http.Get(url)
	if err != nil {
		respondWithServerError(w, err)
		return
	}
	defer response.Body.Close()

	// file not found
	if response.StatusCode != http.StatusOK {
		respondWithUserError(w, "File not found")
		return
	}
	// read response
	download, err := ioutil.ReadAll(response.Body)
	if err != nil {
		respondWithServerError(w, err)
		return
	}
	// log.Printf("size of download from %s: %d", url, len(download))

	// create buffer and scan
	part := ioutil.NopCloser(bytes.NewBuffer(download))
	result := a.scan(part)
	// write result
	fileResult := FileResult{"download", result}
	scanResult = append(scanResult, fileResult)
	logger(url, "download", result)

	respondWithJSON(w, http.StatusOK, scanResult)
}

// used by kubernetes
// if this fails kubernetes will restart the container
func (a *App) livenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// used by kubernetes
// if this fails kubernetes stops routing traffic to the container
func (a *App) readinessCheck(w http.ResponseWriter, r *http.Request) {
	version, err := a.Clam.Version()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(version.Raw))
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
	respondWithJSONError(w, http.StatusBadRequest, "Bad Request", description)
}

func respondWithServerError(w http.ResponseWriter, error error) {
	respondWithJSONError(w, http.StatusInternalServerError, "Internal Server Error", error.Error())
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
