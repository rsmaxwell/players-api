package main

import (
	"fmt"
	"os"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pkg          = debug.NewPackage("main")
	functionMain = debug.NewFunction(pkg, "main")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-createdb", "https://server.rsmaxwell.co.uk/archiva")
}

// http://go-database-sql.org/retrieving.html
func main() {
	f := functionMain
	f.Infof("Players CreateDB: Version: %s", basic.Version())

	// Read configuration and connect to the database
	db, _, err := config.Setup()
	if err != nil {
		message := "Error setting up"
		f.Errorf(message)
		f.DumpError(err, message)
		os.Exit(1)
	}
	defer db.Close()

	// Disconnect all users from the database, except ourselves
	sqlStatement := `
	SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.datname = 'players'
	  AND pid <> pg_backend_pid()`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not disconnect other users from database"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Drop the database
	sqlStatement = `DROP DATABASE players`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not drop database"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	fmt.Println("Successfully deleted the database.")
}
