package main

import (
	"fmt"
	"os"

	"github.com/rsmaxwell/players-api/internal/backup"
	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pkg                 = debug.NewPackage("main")
	functionMain        = debug.NewFunction(pkg, "main")
	functionListPeople  = debug.NewFunction(pkg, "listPeople")
	functionListCourts  = debug.NewFunction(pkg, "listCourts")
	functionListPlaying = debug.NewFunction(pkg, "listPlaying")
	functionListWaiting = debug.NewFunction(pkg, "listWaiting")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-createdb", "https://server.rsmaxwell.co.uk/archiva")
}

// http://go-database-sql.org/retrieving.html
func main() {
	f := functionMain
	f.Infof("Players List: Version: %s", basic.Version())

	// Read configuration and connect to the database
	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	err = listPeople(db)
	if err != nil {
		message := "Could not list the people"
		f.Errorf(message)
		os.Exit(1)
	}

	err = listCourts(db)
	if err != nil {
		message := "Could not list the courts"
		f.Errorf(message)
		os.Exit(1)
	}

	err = listPlaying(db)
	if err != nil {
		message := "Could not list the playing"
		f.Errorf(message)
		os.Exit(1)
	}

	err = listWaiting(db)
	if err != nil {
		message := "Could not list the waiting"
		f.Errorf(message)
		os.Exit(1)
	}

	fmt.Printf("Successfully listed database: %s\n", c.Database.DatabaseName)
}

func listPeople(db *sql.DB) error {
	f := functionListPeople

	// Query all the records in the people table

	fields := "id, firstname, lastname, knownas, email, phone, hash, status"
	sqlStatement := "SELECT " + fields + " FROM " + model.PersonTable
	f.DebugVerbose("%s", sqlStatement)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select from the table: " + model.PersonTable
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		os.Exit(1)
	}
	defer rows.Close()

	f.Infof("")
	f.Infof("---[ people ]------------------------")
	var np model.NullPerson
	for rows.Next() {
		err := rows.Scan(&np.ID, &np.FirstName, &np.LastName, &np.Knownas, &np.Email, &np.Phone, &np.Hash, &np.Status)
		if err != nil {
			message := "could not scan the person record"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		f.Infof("id: %d", np.ID)

		if np.FirstName.Valid {
			f.Infof("firstName:%s", np.FirstName.String)
		}

		if np.LastName.Valid {
			f.Infof("lastName:%s", np.LastName.String)
		}

		if np.Knownas.Valid {
			f.Infof("knownas:%s", np.Knownas.String)
		}

		if np.Email.Valid {
			f.Infof("email:%s", np.Email.String)
		}

		if np.Phone.Valid {
			f.Infof("phone:%s", np.Phone.String)
		}

		if np.Hash.Valid {
			f.Infof("hash:%s", np.Hash.String)
		}

		if np.Status.Valid {
			f.Infof("status:%s", np.Status.String)
		}

		f.Infof("-------------------------------------")
	}
	err = rows.Err()
	if err != nil {
		message := "error iterating through the query results"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}
	f.Infof("")
	return nil
}

func listCourts(db *sql.DB) error {
	f := functionListCourts

	// Query all the records in the courts table
	sqlStatement := "SELECT * FROM " + model.CourtTable
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select from " + model.CourtTable
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	f.Infof("")
	f.Infof("---[ courts ]------------------------")
	var c model.NullCourt
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			message := "error scanning the results"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		f.Infof("id: %d", c.ID)

		if c.Name.Valid {
			f.Infof("name:%s", c.Name.String)
		}
		f.Infof("-------------------------------------")
	}
	err = rows.Err()
	if err != nil {
		message := "error iterating through the query results"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}
	f.Infof("")
	return nil
}

func listPlaying(db *sql.DB) error {
	f := functionListPlaying

	// Query all the records in the playing table
	sqlStatement := `SELECT * FROM playing`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all playing"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	f.Infof("")
	f.Infof("---[ playing ]-----------------------")
	var p backup.Play
	for rows.Next() {
		err := rows.Scan(&p.Person, &p.Court)
		if err != nil {
			message := "error scanning the results"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		f.Infof("person: %d", p.Person)
		f.Infof("court:  %d", p.Court)
		f.Infof("-------------------------------------")
	}
	err = rows.Err()
	if err != nil {
		message := "error iterating through the query results"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}
	f.Infof("")
	return nil
}

func listWaiting(db *sql.DB) error {
	f := functionListWaiting

	// Query all the records in the waiting table
	sqlStatement := `SELECT * FROM waiting`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all waiting"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	f.Infof("")
	f.Infof("---[ waiting ]-----------------------")
	var p int
	var nt sql.NullTime
	for rows.Next() {
		err := rows.Scan(&p, &nt)
		if err != nil {
			message := "error scanning the results"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		f.Infof("person: %d", p)

		if nt.Valid {
			f.Infof("start:  %s", nt.Time.Format("2006-01-02 15:04:05.999999"))
		} else {
			f.Infof("start:  (nil)")
		}
		f.Infof("-------------------------------------")
	}
	err = rows.Err()
	if err != nil {
		message := "error iterating through the query results"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}
	f.Infof("")
	return nil
}
