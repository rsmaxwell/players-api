package main

import (
	"fmt"
	"os"

	"github.com/rsmaxwell/players-api/internal/model"

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
	f.Infof("Players CreateTables: Version: %s", basic.Version())

	AdminFirstName, ok := os.LookupEnv("PLAYERS_ADMIN_FIRST_NAME")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_FIRST_NAME not set")
		os.Exit(1)
	}

	AdminLastName, ok := os.LookupEnv("PLAYERS_ADMIN_LAST_NAME")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_LAST_NAME not set")
		os.Exit(1)
	}

	AdminDisplayName, ok := os.LookupEnv("PLAYERS_ADMIN_DISPLAY_NAME")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_DISPLAY_NAME not set")
		os.Exit(1)
	}

	AdminUserName, ok := os.LookupEnv("PLAYERS_ADMIN_USERNAME")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_USERNAME not set")
		os.Exit(1)
	}

	AdminEmail, ok := os.LookupEnv("PLAYERS_ADMIN_EMAIL")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_EMAIL not set")
		os.Exit(1)
	}

	AdminPhone, ok := os.LookupEnv("PLAYERS_ADMIN_PHONE")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_PHONE not set")
		os.Exit(1)
	}

	AdminPassword, ok := os.LookupEnv("PLAYERS_ADMIN_PASSWORD")
	if !ok {
		f.Errorf("PLAYERS_ADMIN_PASSWORD not set")
		os.Exit(1)
	}

	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	// Drop the playing table
	sqlStatement := `DROP TABLE ` + model.PlayingTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not drop playing"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
	}

	// Drop the waiting table
	sqlStatement = `DROP TABLE ` + model.WaitingTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not drop waiting"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
	}

	// Drop the person table
	sqlStatement = "DROP TABLE " + model.PersonTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not drop person"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
	}

	// Drop the court table
	sqlStatement = "DROP TABLE " + model.CourtTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not drop court"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
	}

	// Create the people table
	sqlStatement = `
		CREATE TABLE ` + model.PersonTable + ` (
			id SERIAL PRIMARY KEY,
			username VARCHAR(64) NOT NULL UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			displayname VARCHAR(32) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			phone VARCHAR(32) NOT NULL UNIQUE,
			hash VARCHAR(255) NOT NULL,	
			status VARCHAR(32) NOT NULL
		 )`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create person table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Create the person_email index
	sqlStatement = "CREATE INDEX person_email ON " + model.PersonTable + " ( email )"
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create person_email index"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Create the person_username index
	sqlStatement = "CREATE INDEX person_username ON " + model.PersonTable + " ( username )"
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create person_username index"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Create the court table
	sqlStatement = `
		CREATE TABLE ` + model.CourtTable + ` (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255)
		 )`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create court table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Create the playing table
	sqlStatement = `
		CREATE TABLE ` + model.PlayingTable + ` (
			person INT PRIMARY KEY,
			court  INT,

			CONSTRAINT person FOREIGN KEY(person) REFERENCES person(id),
			CONSTRAINT court FOREIGN KEY(court)  REFERENCES court(id)
		 )`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create playing table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	// Create the waiting table
	sqlStatement = `
		CREATE TABLE ` + model.WaitingTable + ` (
			person INT PRIMARY KEY,
			start  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT person FOREIGN KEY(person) REFERENCES person(id)
		 )`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not create waiting table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}

	fmt.Printf("Successfully created Tables in the database: %s\n", c.Database.DatabaseName)

	peopleData := []model.Registration{
		{FirstName: AdminFirstName, LastName: AdminLastName, DisplayName: AdminDisplayName, UserName: AdminUserName, Email: AdminEmail, Phone: AdminPhone, Password: AdminPassword},
	}

	peopleIDs := make(map[int]int)
	for i, r := range peopleData {

		p, err := r.ToPerson()
		if err != nil {
			message := "Could not register person"
			f.Errorf(message)
			f.DumpError(err, message)
			os.Exit(1)
		}

		p.Status = model.StatusAdmin

		err = p.SavePerson(db)
		if err != nil {
			message := fmt.Sprintf("Could not save person: firstName: %s, lastname: %s, username: %s, email: %s", p.FirstName, p.LastName, p.UserName, p.Email)
			f.Errorf(message)
			f.DumpError(err, message)
			os.Exit(1)
		}

		peopleIDs[i] = p.ID
	}
}
