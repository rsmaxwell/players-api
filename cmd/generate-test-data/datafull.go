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
		role      string
		player    bool
	}{
		{id: "007", password: "topsecret", firstname: "James", lastname: "Bond", email: "james@mi6.co.uk", role: model.RoleAdmin, player: true},
		{id: "bob", password: "qwerty", firstname: "Robert", lastname: "Bruce", email: "bob@aol.com", role: model.RoleNormal, player: true},
		{id: "alice", password: "wonder", firstname: "Alice", lastname: "Wonderland", email: "alice@abc.com", role: model.RoleNormal, player: true},
		{id: "jill", password: "password", firstname: "Jill", lastname: "Cooper", email: "jill@def.com", role: model.RoleNormal, player: true},
		{id: "david", password: "magic", firstname: "David", lastname: "Copperfield", email: "david@ghi.com", role: model.RoleNormal, player: true},
		{id: "mary", password: "queen", firstname: "Mary", lastname: "Gray", email: "mary@jkl.com", role: model.RoleNormal, player: true},
		{id: "john", password: "king", firstname: "John", lastname: "King", email: "john@mno.com", role: model.RoleNormal, player: true},
		{id: "judith", password: "bean", firstname: "Judith", lastname: "Green", email: "james@mi6.co.uk", role: model.RoleNormal, player: true},
		{id: "paul", password: "ruler", firstname: "Paul", lastname: "Straight", email: "paul@stu.com", role: model.RoleNormal, player: true},
		{id: "nigel", password: "changeme", firstname: "suspended", lastname: "Nonplayer", email: "nigel@vwx.com", role: model.RoleSuspended, player: false},
		{id: "jeremy", password: "danger", firstname: "suspended", lastname: "Player", email: "jeremy@vwx.com", role: model.RoleSuspended, player: true},
		{id: "joanna", password: "bright", firstname: "Nonplayer", lastname: "Brown", email: "joanna@yza.com", role: model.RoleNormal, player: false},
	}

	for _, i := range datapeople {
		err = Register(i.id, i.password, i.firstname, i.lastname, i.email)
		if err != nil {
			return err
		}
	}

	// Set the security role of the people
	for _, i := range datapeople {
		err = model.UpdatePersonRole(i.id, i.role)
		if err != nil {
			return err
		}
	}

	// Set the 'player' field of the people
	for _, i := range datapeople {
		err = model.UpdatePersonPlayer(i.id, i.player)
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
