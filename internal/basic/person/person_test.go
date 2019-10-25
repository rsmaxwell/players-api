package person

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/rsmaxwell/players-api/internal/codeerror"

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

func TestNewInfoJunkPerson(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("topsecret"), bcrypt.DefaultCost)
	require.Nil(t, err, "err should be nothing")

	list1, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = New("John", "Buutler", "butler@ghj.co.uk", hashedPassword, false).Save("butler")
	require.Nil(t, err, "err should be nothing")

	list2, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1)+1, len(list2), "List of people not updated correctly")
}

func TestSavePerson(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	list1, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = New("fred", "smith", junkEmail, junkHashedPassword, false).Save("fred")
	require.Nil(t, err, "err should be nothing")

	list2, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1)+1, len(list2))
}

func TestPerson(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Get the initial number of people
	list1, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")

	// Add a couple of people
	datapeople := []struct {
		id             string
		firstname      string
		lastname       string
		email          string
		hashedpassword []byte
		role           string
		player         bool
	}{
		{id: "1", firstname: "Fred", lastname: "xxx", email: junkEmail, hashedpassword: junkHashedPassword, role: "admin", player: false},
		{id: "2", firstname: "Bloggs", lastname: "xxx", email: junkEmail, hashedpassword: junkHashedPassword, role: "suspended", player: false},
		{id: "3", firstname: "Jane", lastname: "xxx", email: junkEmail, hashedpassword: junkHashedPassword, role: "suspended", player: false},
		{id: "4", firstname: "Alice", lastname: "xxx", email: junkEmail, hashedpassword: junkHashedPassword, role: "suspended", player: false},
		{id: "5", firstname: "Bob", lastname: "xxx", email: junkEmail, hashedpassword: junkHashedPassword, role: "suspended", player: false},
	}

	// Add a couple of people
	for _, i := range datapeople {
		err := New(i.firstname, i.lastname, i.email, i.hashedpassword, i.player).Save(i.id)
		require.Nil(t, err, "err should be nothing")
	}

	// Check the expected number of People have been created
	list2, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1)+len(datapeople), len(list2), fmt.Sprintf("Unexpected number of people created. expected:%d actual:%d", len(list1)+len(datapeople), len(list2)))

	// Check the expected People have been created
	for _, i := range datapeople {
		p, err := Load(i.id)
		require.Nil(t, err, "err should be nothing")

		require.Equal(t, p.FirstName, i.firstname, fmt.Sprintf("Person[%s] not updated correctly: 'Firstname': expected %s, got: %s", i.id, i.firstname, p.FirstName))
		require.Equal(t, p.LastName, i.lastname, fmt.Sprintf("Person[%s] not updated correctly: 'LastName': expected %s, got: %s", i.id, i.lastname, p.LastName))
		require.Equal(t, p.Email, i.email, fmt.Sprintf("Person[%s] not updated correctly: 'Email': expected %s, got: %s", i.id, i.email, p.Email))
		require.Equal(t, p.Role, i.role, fmt.Sprintf("Person[%s] not updated correctly: 'Role': expected %s, got: %s", i.id, i.role, p.Role))
		require.Equal(t, p.Player, i.player, fmt.Sprintf("Person[%s] not updated correctly: 'Player': expected %t, got: %t", i.id, i.player, p.Player))
	}

	// Delete the people we created
	for _, i := range datapeople {
		err := Remove(i.id)
		require.Nil(t, err, "err should be nothing")
	}

	// Check the final number of people
	list3, err := List(AllRoles)
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1), len(list3), fmt.Sprintf("Unexpected number of people. Expected:%d, actual:%d", len(list1), len(list3)))
}

func TestDeletePersonWithDuffID(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Attempt to delete a person using a duff ID
	err := Remove("junk")
	if err == nil {
		require.Fail(t, fmt.Sprintf("Expected an error. actually got: [%v].", err))
	} else {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusNotFound {
				require.Fail(t, fmt.Sprintf("Unexpected error: [%v]", err))
			}
		} else {
			require.Fail(t, fmt.Sprintf("Unexpected error: [%v]", err))
		}
	}
}

func TestListPeopleWithDuffPlayerFile(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	filename, err := makeFilename("x")
	require.Nil(t, err, "err should be nothing")

	// Create a new person file with junk contents
	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Attempt to use the 'junk' person!
	_, err = List([]string{"regular", "admin", "suspended"})
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				require.Fail(t, fmt.Sprintf("Unexpected error type: expected: %d, Actual: %d", http.StatusNotFound, cerr.Code()))
			}
		} else {
			require.Fail(t, fmt.Sprintf("%s", err))
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}

func TestDetailsWithDuffPersonFile(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Create a new person file with junk contents
	filename, err := makeFilename("0")
	require.Nil(t, err, "err should be nothing")

	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Check that List returns an error
	expected := http.StatusInternalServerError
	_, err = Load("0")
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != expected {
				require.Fail(t, fmt.Sprintf("Unexpected error type: expected: %d, Actual: %d", expected, cerr.Code()))
			}
		} else {
			require.Fail(t, fmt.Sprintf("%s", err))
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}

// removeContents empties the contents of a directory
func removeContents(dirname string) error {
	children, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, d := range children {
		err = os.RemoveAll(path.Join([]string{dirname, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupTestcase function
func setupTestcase(t *testing.T) func(t *testing.T) {

	err := removeContents(personListDir)
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}
