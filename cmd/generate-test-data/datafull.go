package main

import (
	"github.com/rsmaxwell/players-api/internal/model"
)

// createTestdataFull function
func createTestdataFull() error {

	err := clearModel()
	if err != nil {
		return err
	}

	datapeople := []struct {
		id        string
		password  string
		firstname string
		lastname  string
		email     string
	}{
		{id: "007", password: "topsecret", firstname: "James", lastname: "Bond", email: "james@mi6.co.uk"},
		{id: "bob", password: "qwerty", firstname: "Robert", lastname: "Bruce", email: "bob@aol.com"},
		{id: "alice", password: "wonder", firstname: "Alice", lastname: "Wonderland", email: "alice@abc.com"},
		{id: "jill", password: "password", firstname: "Jill", lastname: "Cooper", email: "jill@def.com"},
		{id: "david", password: "magic", firstname: "David", lastname: "Copperfield", email: "david@ghi.com"},
		{id: "mary", password: "queen", firstname: "Mary", lastname: "Gray", email: "mary@jkl.com"},
		{id: "john", password: "king", firstname: "John", lastname: "King", email: "john@mno.com"},
		{id: "judith", password: "bean", firstname: "Judith", lastname: "Green", email: "james@mi6.co.uk"},
		{id: "paul", password: "ruler", firstname: "Paul", lastname: "Straight", email: "paul@stu.com"},
		{id: "nigel", password: "careful", firstname: "Nigel", lastname: "Curver", email: "nigel@vwx.com"},
		{id: "jeremy", password: "changeme", firstname: "Jeremy", lastname: "Black", email: "jeremy@vwx.com"},
		{id: "joanna", password: "bright", firstname: "Joanna", lastname: "Brown", email: "joanna@yza.com"},
	}

	for _, i := range datapeople {
		err = Register(i.id, i.password, i.firstname, i.lastname, i.email)
		if err != nil {
			return err
		}
	}

	// Make all the people a 'player'
	for _, i := range datapeople {
		err = model.UpdatePersonPlayer(i.id, true)
		if err != nil {
			return err
		}
	}

	datacourts := []struct {
		id      string
		players []string
	}{
		{id: "Court 1", players: []string{}},
		{id: "Court 2", players: []string{}},
	}

	for _, i := range datacourts {
		err = CreateCourt(i.id, i.players)
		if err != nil {
			return err
		}
	}

	return nil
}
