package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/basic/version"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/httphandler"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/response"

	"github.com/gorilla/mux"
)

var (
	port               int
	pkg                = debug.NewPackage("main")
	functionMain       = debug.NewFunction(pkg, "main")
	functionMiddleware = debug.NewFunction(pkg, "Middleware")
)

func main() {
	f := functionMain

	f.Infof("Players API: BuildID: %s", version.BuildID())
	f.Verbosef("    BuildDate: %s", version.BuildDate())
	f.Verbosef("    GitCommit: %s", version.GitCommit())
	f.Verbosef("    GitBranch: %s", version.GitBranch())
	f.Verbosef("    GitURL:    %s", version.GitURL())
	var ok bool

	portstring, ok := os.LookupEnv("PORT")
	if !ok {
		portstring = "4201"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		f.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		f.Fatalf(err.Error())
	}

	f.Verbosef("Registering Router and setting Handlers")
	router := mux.NewRouter()
	httphandler.SetupHandlers(router)

	f.Infof("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), Middleware(router))
	if err != nil {
		f.Fatalf(err.Error())
	}
}

// Middleware method
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		f := functionMiddleware

		rw2 := response.New(rw)

		f.DebugRequest(req)
		h.ServeHTTP(rw2, req)
	})
}
