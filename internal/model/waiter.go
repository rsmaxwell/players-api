package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rsmaxwell/players-api/internal/debug"
)

// Waiter type
type Waiter struct {
	Person int       `json:"person"`
	Start  time.Time `json:"start"`
}

// NullWaiter type
type NullWaiter struct {
	Person int
	Start  sql.NullTime
}

const (
	// WaitingTable is the name of the waiting table
	WaitingTable = "waiting"
)

var (
	functionListWaiters          = debug.NewFunction(pkg, "ListWaiters")
	functionListWaitersForPerson = debug.NewFunction(pkg, "ListWaitersForPerson")
)

// ListWaiters returns the list of waiters
func ListWaiters(db *sql.DB) ([]Waiter, error) {
	f := functionListWaiters

	sqlStatement := "SELECT * FROM " + WaitingTable + " ORDER BY start ASC"

	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not get list the waiters"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []Waiter
	for rows.Next() {

		var nw NullWaiter
		err := rows.Scan(&nw.Person, &nw.Start)
		if err != nil {
			message := "Could not scan the waiter"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		var w Waiter
		w.Person = nw.Person

		if nw.Start.Valid {
			w.Start = nw.Start.Time
		}

		list = append(list, w)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the waiters"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// ListWaitersForPerson returns the list of waiters for a person
func ListWaitersForPerson(db *sql.DB, id int) ([]Waiter, error) {
	f := functionListWaitersForPerson

	fields := "person, start"
	sqlStatement := "SELECT " + fields + " FROM " + WaitingTable + " WHERE person=$1"

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		message := "Could not get list the waiters"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []Waiter
	for rows.Next() {

		var nw NullWaiter
		err := rows.Scan(&nw.Person, &nw.Start)
		if err != nil {
			message := "Could not scan the waiter"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		var w Waiter
		w.Person = nw.Person

		if nw.Start.Valid {
			w.Start = nw.Start.Time
		}

		list = append(list, w)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the waiters"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// Dump writes the waiter to a dump file
func (w *Waiter) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(w)
	if err != nil {
		return
	}

	title := fmt.Sprintf("waiter.json")
	d.AddByteArray(title, bytearray)
}
