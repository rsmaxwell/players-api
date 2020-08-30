package model

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/rsmaxwell/players-api/internal/debug"
)

const (
	// GoodFirstName const
	GoodFirstName = "James"

	// GoodLastName const
	GoodLastName = "Bond"

	// GoodDisplayName const
	GoodDisplayName = "007"

	// GoodUserName const
	GoodUserName = "007"

	// GoodEmail const
	GoodEmail = "007@mi6.gov.uk"

	// GoodPhone const
	GoodPhone = "+44 000 000000"

	// GoodPassword const
	GoodPassword = "TopSecret"

	// AnotherFirstName const
	AnotherFirstName = "Alice"

	// AnotherLastName const
	AnotherLastName = "Smith"

	// AnotherDisplayName const
	AnotherDisplayName = "Alice"

	// AnotherUserName const
	AnotherUserName = "Ally"

	// AnotherEmail const
	AnotherEmail = "alice@aol.com"

	// AnotherPhone const
	AnotherPhone = "07856 123456"

	// AnotherPassword const
	AnotherPassword = "darkblue"
)

// Logon type
type Logon struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,len=30"`
}

var (
	pkg = debug.NewPackage("model")

	functionRegister     = debug.NewFunction(pkg, "Register")
	functionMakePlaying  = debug.NewFunction(pkg, "MakePlaying")
	functionMakeWaiting  = debug.NewFunction(pkg, "MakeWaiting")
	functionMakeInactive = debug.NewFunction(pkg, "MakeInactive")
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
	functionDeleteCourt  = debug.NewFunction(pkg, "DeleteCourt")
	functionPopulate     = debug.NewFunction(pkg, "Populate")
)

// Populate adds a new set of standard records
func Populate(db *sql.DB) error {
	f := functionPopulate

	peopleData := []Registration{
		{FirstName: GoodFirstName, LastName: GoodLastName, DisplayName: GoodDisplayName, UserName: GoodUserName, Email: GoodEmail, Phone: GoodPhone, Password: GoodPassword},
		{FirstName: AnotherFirstName, LastName: AnotherLastName, DisplayName: AnotherDisplayName, UserName: AnotherUserName, Email: AnotherEmail, Phone: AnotherPhone, Password: AnotherPassword},
		{FirstName: "Robert", LastName: "Brown", DisplayName: "Bob", UserName: "bob1843", Email: "bob@ntl.co.uk", Phone: "012345 123010", Password: "Browneyes"},
		{FirstName: "Charles", LastName: "Winsor", DisplayName: "Charlie", UserName: "cw8765", Email: "charles@o2.co.uk", Phone: "012345 123011", Password: "hrhcharles"},
		{FirstName: "David", LastName: "Townsend", DisplayName: "Dave", UserName: "dtownsend1970", Email: "david@bt.co.uk", Phone: "012345 123012", Password: "miltonkeynes"},
		{FirstName: "Edward", LastName: "French", DisplayName: "Ed", UserName: "efrench87", Email: "immissadda-1167@yopmail.com", Phone: "012345 123013", Password: "romeroandjuliet"},
		{FirstName: "Hana", LastName: "Johnson", DisplayName: "Han", UserName: "hjohn7654", Email: "uddobareqi-9086@yopmail.com", Phone: "012345 123014", Password: "tabithathecat"},
		{FirstName: "Annette", LastName: "Mack", DisplayName: "Nettie", UserName: "amack456", Email: "benagassuf-0898@yopmail.com", Phone: "012345 123015", Password: "kayleightown"},
		{FirstName: "Karen", LastName: "Curry", DisplayName: "Kara", UserName: "kcurry45", Email: "pyffacisi-2285@yopmail.com", Phone: "012345 123016", Password: "sparkleykeira"},
		{FirstName: "Halima", LastName: "Frazier", DisplayName: "Hal", UserName: "hf1234", Email: "esunnassuppa-5488@yopmail.com", Phone: "012345 123017", Password: "glitterma"},
		{FirstName: "Laila", LastName: "Mcgrath", DisplayName: "La", UserName: "lmcgrath98", Email: "enarula-8425@yopmail.com", Phone: "012345 123018", Password: "tinkerham"},
		{FirstName: "Caroline", LastName: "Clarke", DisplayName: "Carol", UserName: "cc7654", Email: "hossemmibe-4189@yopmail.com", Phone: "012345 123019", Password: "ruificent"},
	}

	peopleIDs := make(map[int]int)
	for i, r := range peopleData {

		p, err := r.ToPerson()
		if err != nil {
			f.Errorf("Could not register person")
			return err
		}

		p.Status = StatusNormal

		err = p.SavePerson(db)
		if err != nil {
			f.Errorf("Could not save person: firstName: %s, lastname: %s, email: %s", p.FirstName, p.LastName, p.Email)
			return err
		}

		err = MakeWaiting(db, p.ID)
		if err != nil {
			f.Errorf("Could not add waiting")
			return err
		}

		peopleIDs[i] = p.ID
	}

	courtData := []struct {
		name string
	}{
		{"A"},
		{"B"},
	}

	courtIDs := make(map[int]int)
	for i, x := range courtData {
		c := NewCourt(x.name)
		err := c.SaveCourt(db)
		if err != nil {
			message := "Could not add court"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}
		courtIDs[i] = c.ID
	}
	return nil
}

// DeletePerson removes a person and associated waiters and playings
func DeletePerson(db *sql.DB, id int) error {
	f := functionDeletePerson

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	// Remove the associated waiters
	sqlStatement := "DELETE FROM " + WaitingTable + " WHERE person=" + strconv.Itoa(id)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not delete waiters"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the associated playing
	sqlStatement = "DELETE FROM " + PlayingTable + " WHERE person=" + strconv.Itoa(id)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete playings"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the Person
	sqlStatement = "DELETE FROM " + PersonTable + " WHERE ID=" + strconv.Itoa(id)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete person"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		message := "Could not commit a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}

// DeleteCourt removes a court and associated playings
func DeleteCourt(db *sql.DB, id int) error {
	f := functionDeleteCourt

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	// Remove the associated playing
	sqlStatement := "DELETE FROM " + PlayingTable + " WHERE court=" + strconv.Itoa(id)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete playings"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the Court
	sqlStatement = "DELETE FROM " + CourtTable + " WHERE ID=" + strconv.Itoa(id)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete court"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		message := "Could not commit a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}

// MakeWaiting moves a person from playing to waiting
func MakeWaiting(db *sql.DB, personID int) error {
	f := functionMakeWaiting

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	// Remove the person from the playing table
	sqlStatement := "DELETE FROM " + PlayingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from playing table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the person from the waiting table
	sqlStatement = "DELETE FROM " + WaitingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from waiting table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Add the person to the waiting table
	now := time.Now()
	fields := "person, start"
	sqlStatement = "INSERT INTO " + WaitingTable + " (" + fields + ") VALUES ($1, $2)"
	_, err = db.ExecContext(ctx, sqlStatement, personID, now)
	if err != nil {
		tx.Rollback()
		message := "Could not insert person into waiting table"
		f.Errorf(message)
		d := f.DumpSQLError(err, message, sqlStatement)

		w := Waiter{Person: personID, Start: now}
		w.Dump(d)

		p := Person{ID: personID}
		p.LoadPerson(db)
		p.Dump(d)

		listOfWaiters, _ := ListWaiters(db)
		d.AddString("NumberOfWaiters.txt", fmt.Sprintf("Number of waiters: %d", len(listOfWaiters)))
		return err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}

// MakePlaying moves a person from playing to waiting
func MakePlaying(db *sql.DB, personID int, courtID int) error {
	f := functionMakePlaying

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	// Remove the person from the playing table
	sqlStatement := "DELETE FROM " + PlayingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from playing table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the person from the waiting table
	sqlStatement = "DELETE FROM " + WaitingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from waiting table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Add the person to the playing table

	fields := "person, court"
	sqlStatement = "INSERT INTO " + PlayingTable + " ( " + fields + " ) VALUES ($1, $2)"
	_, err = db.ExecContext(ctx, sqlStatement, personID, courtID)
	if err != nil {
		tx.Rollback()
		message := "Could not insert person into playing table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}

// MakeInactive removes a player from both the waiting and playing
func MakeInactive(db *sql.DB, personID int) error {
	f := functionMakeInactive

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	// Remove the person from the playing table
	sqlStatement := "DELETE FROM " + PlayingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from playing table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the person from the waiting table
	sqlStatement = "DELETE FROM " + WaitingTable + " WHERE person=" + strconv.Itoa(personID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		tx.Rollback()
		message := "Could not remove person from waiting table"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}
