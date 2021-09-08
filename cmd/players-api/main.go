package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/httphandler"
	"github.com/rsmaxwell/players-api/internal/model"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pkg          = debug.NewPackage("main")
	functionMain = debug.NewFunction(pkg, "main")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-api", "https://server.rsmaxwell.co.uk/archiva")
}

func main() {
	f := functionMain
	ctx := context.Background()

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

	count, err := model.CheckConistencyTx(db, true)
	if err != nil {
		f.Errorf("Error checking consistency")
		os.Exit(1)
	}
	if count != 0 {
		f.Errorf("Inconsistent database: count: %d", count)
		os.Exit(1)
	}

	f.Verbosef("Registering Router and setting Handlers")
	router := mux.NewRouter()
	httphandler.SetupHandlers(router)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	credentials := handlers.AllowCredentials()

	handler := handlers.CORS(headers, methods, origins, credentials)(router)
	handler = httphandler.WithLogging(handler)
	handler = httphandler.AddDatabaseContext(handler, db)
	handler = httphandler.AddRequestContext(handler)
	handler = httphandler.AddConfigContext(handler, c)

	f.Infof("Listening on port: %d", c.Server.Port)
	address := fmt.Sprintf(":%d", c.Server.Port)
	err = http.ListenAndServe(address, handler)
	if err != nil {
		f.Fatalf(ctx, err.Error())
	}
}
