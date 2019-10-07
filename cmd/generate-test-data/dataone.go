package main

import "github.com/rsmaxwell/players-api/internal/basic/person"

// createTestdataOne function
func createTestdataOne() error {

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
		{id: "007", password: "topsecret", firstname: "James", lastname: "Bond", email: "james@mi6.co.uk", role: person.RoleAdmin, player: true},
		{id: "bob", password: "qwerty", firstname: "Robert", lastname: "Bruce", email: "bob@aol.com", role: person.RoleNormal, player: true},
		{id: "alice", password: "wonder", firstname: "Alice", lastname: "Wonderland", email: "alice@abc.com", role: person.RoleSuspended, player: true},
	}

	for _, i := range datapeople {
		err = Register(i.id, i.password, i.firstname, i.lastname, i.email)
		if err != nil {
			return err
		}

		err = person.UpdateRole(i.id, i.role)
		if err != nil {
			return err
		}
	}

	return nil
}
