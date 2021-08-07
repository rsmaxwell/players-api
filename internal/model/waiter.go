package model

import (
	"context"
	"database/sql"
	"encoding/json"
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
	functionGetFirstWaiter       = debug.NewFunction(pkg, "GetFirstWaiter")
	functionRemoveWaiter         = debug.NewFunction(pkg, "RemoveWaiter")
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

// Get first GetFirstWaiter
func GetFirstWaiterContext(db *sql.DB, ctx context.Context) (int, error) {
	f := functionGetFirstWaiter

	fields := "person"
	sqlStatement := "SELECT " + fields + " FROM " + WaitingTable + " LIMIT 1"
	rows, err := db.QueryContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not get the first waiter"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return 0, err
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			message := "Could not scan the first player"
			f.Errorf(message)
			f.DumpError(err, message)
			return 0, err
		}
	}
	err = rows.Err()
	if err != nil {
		message := "Could not get the first players"
		f.Errorf(message)
		f.DumpError(err, message)
		return 0, err
	}

	return id, nil
}

// AddWaiter
func AddWaiter(db *sql.DB, personID int) error {
	return AddWaiterContext(db, context.Background(), personID)
}

// AddWaiter
func AddWaiterContext(db *sql.DB, ctx context.Context, personID int) error {
	f := functionRemoveWaiter

	start := time.Now()

	fields := "person, start"
	values := "$1, $2"
	sqlStatement := "INSERT INTO " + WaitingTable + " (" + fields + ") VALUES (" + values + ")"

	_, err := db.ExecContext(ctx, sqlStatement, personID, start)
	if err != nil {
		message := "Could not insert into " + WaitingTable
		f.Errorf(message)
		d := f.DumpSQLError(err, message, sqlStatement)
		data := struct {
			PersonID int
			Start    time.Time
		}{
			PersonID: personID,
			Start:    start,
		}
		bytes, _ := json.MarshalIndent(data, "", "    ")
		d.AddByteArray("values.json", bytes)
		return err
	}

	return nil
}

// RemoveWaiter
func RemoveWaiter(db *sql.DB, personID int) error {
	return RemoveWaiterContext(db, context.Background(), personID)
}

// RemoveWaiter
func RemoveWaiterContext(db *sql.DB, ctx context.Context, personID int) error {
	f := functionRemoveWaiter

	sqlStatement := "DELETE FROM " + WaitingTable + " WHERE person=$1"
	rows, err := db.QueryContext(ctx, sqlStatement, personID)
	if err != nil {
		message := "Could not delete the waiter"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}
	defer rows.Close()

	return nil
}

// Dump writes the waiter to a dump file
func (w *Waiter) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(w)
	if err != nil {
		return
	}

	title := "waiter.json"
	d.AddByteArray(title, bytearray)
}
