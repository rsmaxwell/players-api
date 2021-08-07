package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionSetup            = debug.NewFunction(pkg, "Setup")
	functionDeleteAllRecords = debug.NewFunction(pkg, "DeleteAllRecords")
	functionFillCourt        = debug.NewFunction(pkg, "FillCourt")
	functionClearCourt       = debug.NewFunction(pkg, "ClearCourt")
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

// FillCourt
func FillCourt(db *sql.DB, courtID int) ([]Position, error) {
	f := functionFillCourt

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	players, err := ListPlayersForCourtContext(db, ctx, courtID)
	if err != nil {
		tx.Rollback()
		message := "Could not list players"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	data, _ := json.MarshalIndent(players, "", "    ")
	f.Infof("players: %s", string(data))

	mapOfPlayers := make(map[int]*Player)
	for _, player := range players {
		p := player                        // take a copy of the object ...
		mapOfPlayers[player.Position] = &p // ... so their addresses are actually different!
	}

	changes := 0
	positions := make([]Position, 0)
	for index := 0; index < NumberOfCourtPositions; index++ {

		var ok bool
		var player *Player
		var personID int

		if player, ok = mapOfPlayers[index]; !ok {
			changes++

			personID, err = GetFirstWaiterContext(db, ctx)
			if err != nil {
				tx.Rollback()
				message := "Could not get the first waiter"
				f.Errorf(message)
				f.DumpError(err, message)
				return nil, err
			}
			f.Infof(fmt.Sprintf("Got First Waiter: [%d]", personID))

			err = RemoveWaiterContext(db, ctx, personID)
			if err != nil {
				tx.Rollback()
				message := "Could not remove the waiter"
				f.Errorf(message)
				f.DumpError(err, message)
				return nil, err
			}
			f.Infof(fmt.Sprintf("Removed Waiter: [%d]", personID))

			err = AddPlayerContext(db, ctx, personID, courtID, index)
			if err != nil {
				tx.Rollback()
				message := "Could not add player"
				f.Errorf(message)
				f.DumpError(err, message)
				return nil, err
			}
			p := Player{Person: personID, Court: courtID, Position: index}
			player = &p
			f.Infof(fmt.Sprintf("Added Player: [person: %d, courtID:%d, position:%d]", personID, courtID, index))
		}

		person := FullPerson{ID: player.Person}
		err = person.LoadPerson(db)
		if err != nil {
			tx.Rollback()
			message := "Could not load player"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		var position = Position{Index: player.Position, PersonID: player.Person, DisplayName: person.Knownas}
		positions = append(positions, position)
	}

	if changes > 0 {
		err = tx.Commit()
		if err != nil {
			message := "Could not commit transaction"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}
	} else {
		err = tx.Rollback()
		if err != nil {
			message := "Could not rollback transaction"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}
	}

	return positions, nil
}

// ClearCourt
func ClearCourt(db *sql.DB, courtID int) error {
	f := functionClearCourt

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin transaction"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	players, err := ListPlayersForCourtContext(db, ctx, courtID)
	if err != nil {
		tx.Rollback()
		message := "Could not list players"
		f.Errorf(message)
		f.DumpError(err, message)
		return err
	}

	for _, player := range players {

		err = RemovePlayerContext(db, ctx, player.Person)
		if err != nil {
			tx.Rollback()
			message := "Could not remove player"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}

		err = AddWaiterContext(db, ctx, player.Person)
		if err != nil {
			tx.Rollback()
			message := "Could not add the waiter"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}
	}

	if len(players) > 0 {
		err = tx.Commit()
		if err != nil {
			message := "Could not commit transaction"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}
	} else {
		err = tx.Rollback()
		if err != nil {
			message := "Could not rollback transaction"
			f.Errorf(message)
			f.DumpError(err, message)
			return err
		}
	}

	return nil
}

// EqualIntArray tells whether a and b contain the same elements NOT in-order order
func EqualIntArray(x, y []int) bool {

	if x == nil {
		return y == nil
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
