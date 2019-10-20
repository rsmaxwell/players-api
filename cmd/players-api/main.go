package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/basic/version"
	"github.com/rsmaxwell/players-api/internal/httphandler"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
)

var (
	port int
)

func main() {

	log.Printf("Players Server: %s", version.GetVersion())
	var ok bool

	portstring, ok := os.LookupEnv("PORT")
	if !ok {
		portstring = "4201"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Registering Router and setting Handlers")
	router := mux.NewRouter()
	httphandler.SetupHandlers(router)

	log.Printf("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Success")
}
