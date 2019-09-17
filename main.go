package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/httpHandler"
	"github.com/rsmaxwell/players-api/logger"
)

var (
	port int
)

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	r.HandleFunc("/register",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.Register(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/login",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.Login(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.ListCourts(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.CreateCourt(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.UpdateCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.DeleteCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/person",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.ListPeople(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetPerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.UpdatePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.DeletePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/metrics",
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetMetrics(w, req)
		}).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.NotFound(w, req)
		})
}

func main() {

	logger.Logger.Printf("Players Server: 2018-01-31 13:30")
	var ok bool

	portstring, ok := os.LookupEnv("PORT")
	if !ok {
		portstring = "4201"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Registering Router and setting Handlers")
	router := mux.NewRouter()
	setupHandlers(router)

	logger.Logger.Printf("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Success")
}
