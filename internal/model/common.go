package model

import (
	"database/sql"
	"testing"

	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionSetup            = debug.NewFunction(pkg, "Setup")
	functionDeleteAllRecords = debug.NewFunction(pkg, "DeleteAllRecords")
)

var (
	// MetricsData containing metrics
	MetricsData Metrics
)

// Metrics structure
type Metrics struct {
	StatusCodes map[int]int `json:"statusCodes"`
}

func init() {
	MetricsData = Metrics{}
	MetricsData.StatusCodes = make(map[int]int)
}

// Setup function
func Setup(t *testing.T) (func(t *testing.T), *sql.DB, *config.Config) {
	f := functionSetup

	// Read configuration
	db, c, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		t.FailNow()
	}

	// Delete all the records
	err = DeleteAllRecords(db)
	if err != nil {
		f.Errorf("Error delete all the records")
		t.FailNow()
	}

	// Populate
	err = Populate(db)
	if err != nil {
		f.Errorf("Could not populate the database")
		t.FailNow()
	}

	return func(t *testing.T) {
		db.Close()
	}, db, c
}

// DeleteAllRecords removes all the records in the database
func DeleteAllRecords(db *sql.DB) error {
	f := functionDeleteAllRecords

	sqlStatement := "DELETE FROM " + PlayingTable
	_, err := db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete all from playing"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	sqlStatement = "DELETE FROM " + WaitingTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete all from waiting"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	sqlStatement = "DELETE FROM " + CourtTable
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete all from courts"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	sqlStatement = "DELETE FROM " + PersonTable + " WHERE status != '" + StatusAdmin + "'"
	_, err = db.Exec(sqlStatement)
	if err != nil {
		message := "Could not delete all from people"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return nil
}

// EqualIntArray tells whether a and b contain the same elements NOT in-order order
func EqualIntArray(x, y []int) bool {

	if x == nil {
		if y == nil {
			return true
		}
		return false
	} else if y == nil {
		return false
	}

	if len(x) != len(y) {
		return false
	}

	xMap := make(map[int]int)
	yMap := make(map[int]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}
