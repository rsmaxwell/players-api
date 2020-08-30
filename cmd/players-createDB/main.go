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
	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	// Create the database
	sqlStatement := fmt.Sprintf("CREATE DATABASE %s", c.Database.DatabaseName)
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := fmt.Sprintf("Could not create database: %s", c.Database.DatabaseName)
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	fmt.Printf("Successfully created database: %s\n", c.Database.DatabaseName)
}
