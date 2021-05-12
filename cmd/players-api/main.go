package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/httphandler"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pkg          = debug.NewPackage("main")
	functionInit = debug.NewFunction(pkg, "init")
	functionMain = debug.NewFunction(pkg, "main")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-api", "https://server.rsmaxwell.co.uk/archiva")
}

func main() {
	f := functionMain
	f.Infof("Players API: Version: %s", basic.Version())
	f.Verbosef("    BuildDate: %s", basic.BuildDate())
	f.Verbosef("    GitCommit: %s", basic.GitCommit())
	f.Verbosef("    GitBranch: %s", basic.GitBranch())
	f.Verbosef("    GitURL:    %s", basic.GitURL())

	f.Verbosef("Read configuration and connect to the database")
	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	f.Verbosef("Registering Router and setting Handlers")
	router := mux.NewRouter()
	httphandler.SetupHandlers(router)

	f.Infof("Listening on port: %d", c.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), httphandler.Middleware(router, db))
	if err != nil {
		f.Fatalf(err.Error())
	}

	fmt.Printf("Successfully populated the database: %s\n", c.Database.DatabaseName)
}
