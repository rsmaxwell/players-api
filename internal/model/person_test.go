package model

import (
	"fmt"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/person"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	junkHashedPassword []byte
	junkEmail          string
)

func init() {
	junkHashedPassword = []byte("123456789012345678901234567890123456789012345678901234567890")
	junkEmail = "123@hotmail.com"
}

func TestOverwriteExistingPerson(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("topsecret"), bcrypt.DefaultCost)
	require.Nil(t, err, "err should be nothing")

	list1, err := person.List(person.AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = person.New("John", "Buutler", "butler@ghj.co.uk", hashedPassword, false).Save("007")
	require.Nil(t, err, "err should be nothing")

	list2, err := person.List(person.AllRoles)
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1), len(list2), "List of people not updated correctly")
}

func TestUpdatePerson(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Add a couple of people
	datapeople := []struct {
		id             string
		firstname      string
		lastname       string
		email          string
		hashedpassword []byte
	}{
		{id: "007", firstname: "aaa", lastname: "bbb", email: junkEmail, hashedpassword: junkHashedPassword},
	}

	for _, i := range datapeople {

		p2 := make(map[string]interface{})

		p2["FirstName"] = i.firstname
		p2["LastName"] = i.lastname
		p2["Email"] = i.email
		p2["HashedPassword"] = i.hashedpassword

		err := person.Update("007", p2)
		require.Nil(t, err, "err should be nothing")

		// Check the expected People have been created
		p, err := person.Load(i.id)
		require.Nil(t, err, "err should be nothing")

		require.Equal(t, p.FirstName, i.firstname, fmt.Sprintf("Person[%s] not updated correctly: 'Firstname': expected %s, got: %s", i.id, i.firstname, p.FirstName))
		require.Equal(t, p.LastName, i.lastname, fmt.Sprintf("Person[%s] not updated correctly: 'LastName': expected %s, got: %s", i.id, i.lastname, p.LastName))
		require.Equal(t, p.Email, i.email, fmt.Sprintf("Person[%s] not updated correctly: 'Email': expected %s, got: %s", i.id, i.email, p.Email))
		require.Equal(t, p.HashedPassword, i.hashedpassword, fmt.Sprintf("Person[%s] not updated correctly: 'HashedPassword': expected %s, got: %s", i.id, i.hashedpassword, p.HashedPassword))
	}
}

func TestUpdatePersonPlayer(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	datapeople := []struct {
		id             string
		player         bool
		expectedResult error
	}{
		{id: "007", player: false, expectedResult: nil},
		{id: "alice", player: true, expectedResult: nil},
	}

	for _, i := range datapeople {

		err := person.UpdatePlayer(i.id, i.player)
		require.Nil(t, err, "err should be nothing")

		// Check the expected People have been created
		p, err := person.Load(i.id)
		require.Nil(t, err, "err should be nothing")
		require.Equal(t, i.expectedResult, err, "Unexpected error")
		require.Equal(t, p.Player, i.player, "Unexpected Player: wanted:%s, got:%s", p.Player, i.player)
	}
}

func TestPersonCanLogin(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Add a couple of people
	tests := []struct {
		testName       string
		id             string
		role           string
		player         bool
		expectedResult bool
	}{
		{testName: "admin", id: "007", role: person.RoleAdmin, player: true, expectedResult: true},
		{testName: "suspended non-player", id: "nigel", role: person.RoleSuspended, player: false, expectedResult: false},
		{testName: "suspended player", id: "jeremy", role: person.RoleSuspended, player: true, expectedResult: false},
		{testName: "normal non-player", id: "joanna", role: person.RoleNormal, player: false, expectedResult: true},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			err := person.UpdateRole(test.id, test.role)
			require.Nil(t, err, "err should be nothing")

			err = person.UpdatePlayer(test.id, test.player)
			require.Nil(t, err, "err should be nothing")

			ok := person.CanLogin(test.id)
			require.Equal(t, ok, test.expectedResult, "Unexpected error")
		})
	}
}
