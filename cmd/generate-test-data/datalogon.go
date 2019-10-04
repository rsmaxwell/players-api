package main

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
)

// createTestdataLoggedon function
func createTestdataLoggedon() error {

	err := model.ClearModel()
	if err != nil {
		return err
	}

	var (
		myUserID   = "007"
		myPassword = "topsecret"
	)

	datapeople := []struct {
		id        string
		password  string
		firstname string
		lastname  string
		email     string
	}{
		{id: "007", password: "topsecret", firstname: "James", lastname: "Bond", email: "james@mi6.co.uk"},
	}

	for _, i := range datapeople {
		err = Register(i.id, i.password, i.firstname, i.lastname, i.email)
		if err != nil {
			return err
		}
	}

	myToken, err := Login(myUserID, myPassword)
	if err != nil {
		return err
	}

	mySession := session.LookupToken(myToken)
	if mySession == nil {
		return fmt.Errorf("Failed to lookup token")
	}

	return nil
}
