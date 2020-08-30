package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// Court type
type Court struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required,min=3,max=20"`
}

// NullCourt type
type NullCourt struct {
	ID   int
	Name sql.NullString
}

const (
	// CourtTable is the name of the court table
	CourtTable = "court"
)

var (
	functionUpdateCourt      = debug.NewFunction(pkg, "UpdateCourt")
	functionSaveCourt        = debug.NewFunction(pkg, "SaveCourt")
	functionListCourts       = debug.NewFunction(pkg, "ListCourts")
	functionCourtExists      = debug.NewFunction(pkg, "CourtExists")
	functionLoadCourt        = debug.NewFunction(pkg, "LoadCourt")
	functionDeleteCourtBasic = debug.NewFunction(pkg, "DeleteCourtBasic")
)

// NewCourt initialises a Court object
func NewCourt(name string) *Court {
	c := new(Court)
	c.Name = name
	return c
}

// SaveCourt writes a new Court to disk and returns the generated id
func (c *Court) SaveCourt(db *sql.DB) error {
	f := functionSaveCourt

	fields := ""
	values := ""
	separator := ""

	fields = fields + separator + "name"
	values = values + separator + basic.Quote(c.Name)
	separator = ", "

	sqlStatement := "INSERT INTO " + CourtTable + " (" + fields + ") VALUES (" + values + ") RETURNING id"
	err := db.QueryRow(sqlStatement).Scan(&c.ID)
	if err != nil {
		message := "Could not insert into " + CourtTable
		f.Errorf(message)
		d := f.DumpSQLError(err, message, sqlStatement)
		c.Dump(d)
		return err
	}

	return nil
}

// UpdateCourt method
func (c *Court) UpdateCourt(db *sql.DB) error {
	f := functionUpdateCourt

	items := ""
	separator := ""

	items = items + separator + "name=" + basic.Quote(c.Name)
	separator = ", "

	sqlStatement := "UPDATE " + CourtTable + " SET " + items + " WHERE id=" + strconv.Itoa(c.ID)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		message := "Could not update court"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return err
}

// LoadCourt returns the Court with the given ID
func (c *Court) LoadCourt(db *sql.DB) error {
	f := functionLoadCourt

	// Query the court
	sqlStatement := "SELECT * FROM " + CourtTable + " WHERE ID=" + strconv.Itoa(c.ID)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all people"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++

		var nc NullCourt
		err := rows.Scan(&nc.ID, &nc.Name)
		if err != nil {
			message := "Could not scan the court"
			f.Errorf(message)
			f.DumpError(err, message)
		}

		if nc.Name.Valid {
			c.Name = nc.Name.String
		}
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the courts"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	if count == 0 {
		return codeerror.NewNotFound(fmt.Sprintf("Court id %d not found", c.ID))
	} else if count > 1 {
		message := fmt.Sprintf("Found %d courts with id %d", count, c.ID)
		err := codeerror.NewInternalServerError(message)
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	return nil
}

// DeleteCourtBasic the court with the given ID
func (c *Court) DeleteCourtBasic(db *sql.DB) error {
	f := functionDeleteCourtBasic

	sqlStatement := "DELETE FROM " + CourtTable + " WHERE ID=" + strconv.Itoa(c.ID)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete all from " + CourtTable
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return nil
}

// ListCourts returns a list of the court IDs
func ListCourts(db *sql.DB) ([]int, error) {
	f := functionListCourts

	// Query the court
	sqlStatement := "SELECT id FROM " + CourtTable
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all from " + CourtTable
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []int
	for rows.Next() {

		var nc NullCourt
		err := rows.Scan(&nc.ID)
		if err != nil {
			message := "Could not scan the court"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		list = append(list, nc.ID)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list all from " + CourtTable
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// CourtExists returns 'true' if the court exists
func (c *Court) CourtExists(db *sql.DB) (bool, error) {
	f := functionCourtExists

	// Query the court
	sqlStatement := "SELECT * FROM " + CourtTable + " WHERE id=$1"
	rows, err := db.Query(sqlStatement, c.ID)
	if err != nil {
		message := "Could not select courts"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return false, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the courts"
		f.Errorf(message)
		f.DumpError(err, message)
		return false, err
	}

	if count == 0 {
		return false, nil
	} else if count > 1 {
		message := "Found " + string(count) + " courts with id " + string(c.ID)
		f.Errorf(message)
		f.DumpError(err, message)
		return true, errors.New(message)
	}

	return true, nil
}

// Dump writes the person to a dump file
func (c *Court) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(c)
	if err != nil {
		return
	}

	title := fmt.Sprintf("court.%d.json", c.ID)
	d.AddByteArray(title, bytearray)
}
