package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/httphandler"
)

var (
	port int
)

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	r.HandleFunc("/register",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.Register(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/login",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.Login(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.ListCourts(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.GetCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.CreateCourt(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.UpdateCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.DeleteCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/person",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.ListPeople(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.GetPerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.UpdatePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.DeletePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/metrics",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.GetMetrics(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/queue",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.GetQueue(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/move",
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.PostMove(w, req)
		}).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			httphandler.NotFound(w, req)
		})
}

func main() {

	log.Printf("Players Server: 2018-01-31 13:30")
	var ok bool

	portstring, ok := os.LookupEnv("PORT")
	if !ok {
		portstring = "4201"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Registering Router and setting Handlers")
	router := mux.NewRouter()
	setupHandlers(router)

	log.Printf("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Success")
}
