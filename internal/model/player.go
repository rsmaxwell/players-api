package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rsmaxwell/players-api/internal/debug"
)

// Player type
type Player struct {
	Person int `json:"person"`
	Court  int `json:"court"`
}

// NullPlayer type
type NullPlayer struct {
	Person int
	Court  sql.NullInt32
}

const (
	// PlayingTable is the name of the peplayingrson table
	PlayingTable = "playing"
)

var (
	functionListPlayers          = debug.NewFunction(pkg, "ListPlayers")
	functionListPlayersForPerson = debug.NewFunction(pkg, "ListPlayersForPerson")
)

// ListPlayers returns the list of players
func ListPlayers(db *sql.DB) ([]Player, error) {
	f := functionListPlayers

	sqlStatement := "SELECT * FROM " + PlayingTable
	f.Infof(sqlStatement)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		message := "Could not get list the players"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []Player
	for rows.Next() {

		var np NullPlayer
		err := rows.Scan(&np.Person, &np.Court)
		if err != nil {
			message := "Could not scan the player"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		var p Player
		p.Person = np.Person

		if np.Court.Valid {
			p.Court = int(np.Court.Int32)
		}

		list = append(list, p)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the players"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// ListPlayersForPerson returns the list of players for a person
func ListPlayersForPerson(db *sql.DB, id int) ([]Player, error) {
	f := functionListPlayersForPerson

	fields := "person, court"
	sqlStatement := "SELECT " + fields + " FROM " + PlayingTable + " WHERE person=$1"

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		message := "Could not get list the players"
		f.Errorf(message)
		f.DumpSQLError(err, message, sqlStatement)
		return nil, err
	}
	defer rows.Close()

	var list []Player
	for rows.Next() {

		var np NullPlayer
		err := rows.Scan(&np.Person, &np.Court)
		if err != nil {
			message := "Could not scan the player"
			f.Errorf(message)
			f.DumpError(err, message)
			return nil, err
		}

		var p Player
		p.Person = np.Person

		if np.Court.Valid {
			p.Court = int(np.Court.Int32)
		}

		list = append(list, p)
	}
	err = rows.Err()
	if err != nil {
		message := "Could not list the players"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}

	return list, nil
}

// Dump writes the player to a dump file
func (p *Player) Dump(d *debug.Dump) {

	bytearray, err := json.Marshal(p)
	if err != nil {
		return
	}

	title := fmt.Sprintf("player.json")
	d.AddByteArray(title, bytearray)
}
