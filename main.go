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
	port    int
	baseURL string
)

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	r.HandleFunc(baseURL+"/register", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.Register(w, req)
		})).Methods(http.MethodPost)

	r.HandleFunc(baseURL+"/login", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.Login(w, req)
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/court", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetAllCourts(w, req)
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetCourt(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/court", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.CreateCourt(w, req)
		})).Methods(http.MethodPost)

	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.UpdateCourt(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodPut)

	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.DeleteCourt(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	r.HandleFunc(baseURL+"/person", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetAllPeople(w, req)
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetPerson(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.UpdatePerson(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodPut)

	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.DeletePerson(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	r.HandleFunc(baseURL+"/metrics", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.GetMetrics(w, req)
		})).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			httpHandler.NotFound(w, req)
		}))
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
