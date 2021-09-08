package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pkg              = debug.NewPackage("main")
	functionMain     = debug.NewFunction(pkg, "main")
	functionQueryRow = debug.NewFunction(pkg, "queryRow")
	functionExec     = debug.NewFunction(pkg, "exec")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-createdb", "https://server.rsmaxwell.co.uk/archiva")
}

// http://go-database-sql.org/retrieving.html
func main() {
	f := functionMain
	ctx := context.Background()

	f.Infof("Players Populate: Version: %s", basic.Version())

	// Read configuration and connect to the database
	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	// Delete all the records
	err = model.DeleteAllRecords(ctx, db)
	if err != nil {
		message := "Error delete all the records"
		f.Errorf(message)
		os.Exit(1)
	}

	fmt.Printf("Successfully delete all the records in the database: %s\n", c.Database.DatabaseName)
}
