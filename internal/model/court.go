package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// Position type
type Position struct {
	Index       int    `json:"index"`
	PersonID    int    `json:"personid"`
	DisplayName string `json:"displayname"`
}

// Court type
type Court struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name" validate:"required,min=3,max=20"`
	Positions []Position `json:"positions" db:"positions"`
}

// NullCourt type
type NullCourt struct {
	ID   int
	Name sql.NullString
}

const (
	CourtTable             = "court"
	NumberOfCourtPositions = 4
)

var (
	functionUpdateCourt        = debug.NewFunction(pkg, "UpdateCourt")
	functionSaveCourt          = debug.NewFunction(pkg, "SaveCourt")
	functionListCourts         = debug.NewFunction(pkg, "ListCourts")
	functionLoadCourt          = debug.NewFunction(pkg, "LoadCourt")
	functionDeleteCourt        = debug.NewFunction(pkg, "DeleteCourt")
	functionDeleteCourtContext = debug.NewFunction(pkg, "DeleteCourtContext")
)

// SaveCourt writes a new Court to disk and returns the generated id
func (c *Court) SaveCourt(db *sql.DB) error {
	f := functionSaveCourt

	fields := "name"
	values := basic.Quote(c.Name)

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

	items := "name=" + basic.Quote(c.Name)
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
	return c.LoadCourtContext(db, context.Background())
}

// LoadCourt returns the Court with the given ID
func (c *Court) LoadCourtContext(db *sql.DB, ctx context.Context) error {
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

// DeleteCourt removes a court and associated playings
func (c *Court) DeleteCourt(db *sql.DB) error {
	f := functionDeleteCourt

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	err = DeleteCourtContext(db, ctx, c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit the transaction"
		f.Errorf(message)
		f.DumpError(err, message)
	}

	return nil
}

func DeleteCourtContext(db *sql.DB, ctx context.Context, courtID int) error {
	f := functionDeleteCourtContext

	players, err := ListPlayersForCourtContext(db, ctx, courtID)
	if err != nil {
		message := "Could not delete playings"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	for _, player := range players {
		err = MakePlayerWaitContext(db, ctx, player.Person)
		if err != nil {
			message := "Could not make player wait"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}
	}

	// Remove the associated playing
	sqlStatement := "DELETE FROM " + PlayingTable + " WHERE court=" + strconv.Itoa(courtID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete playings"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	// Remove the Court
	sqlStatement = "DELETE FROM " + CourtTable + " WHERE ID=" + strconv.Itoa(courtID)
	_, err = db.ExecContext(ctx, sqlStatement)
	if err != nil {
		message := "Could not delete court"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return err
	}

	return nil
}

// ListCourts returns a list of the court IDs
func ListCourts(db *sql.DB) ([]Court, error) {
	f := functionListCourts

	// Query the courts
	returnedFields := []string{`id`, `name`}
	sqlStatement := `SELECT ` + strings.Join(returnedFields, `, `) + ` FROM ` + CourtTable
	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not select all from " + CourtTable
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []Court
	for rows.Next() {

		court := Court{}
		court.Positions = make([]Position, 0)

		err := rows.Scan(&court.ID, &court.Name)
		if err != nil {
			message := "Could not scan the court"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		players, err := ListPlayersForCourt(db, court.ID)
		if err != nil {
			message := "Could not list the players on this court"
			f.Errorf(message)
			d := f.DumpError(err, message)

			data, _ := json.MarshalIndent(court, "", "    ")
			d.AddByteArray("court.json", data)

			return nil, err
		}

		for _, player := range players {

			person := FullPerson{ID: player.Person}
			err := person.LoadPerson(db)
			if err != nil {
				message := fmt.Sprintf("Could not load the player [%d]", player.Person)
				f.Errorf(message)
				d := f.DumpError(err, message)
				d.AddObject("court.json", court)
				d.AddObject("player.json", player)
				return nil, err
			}
			position := Position{Index: player.Position, PersonID: player.Person, DisplayName: person.Knownas}
			court.Positions = append(court.Positions, position)
		}

		list = append(list, court)
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

// Dump writes the person to a dump file
func (c *Court) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(c)
	if err != nil {
		return
	}

	title := fmt.Sprintf("court.%d.json", c.ID)
	d.AddByteArray(title, bytearray)
}
