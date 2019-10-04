package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/stretchr/testify/assert"
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

	list1, err := ListPeople(AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = NewPerson("John", "Buutler", "butler@ghj.co.uk", hashedPassword, false).Save("007")
	require.Nil(t, err, "err should be nothing")

	list2, err := ListPeople([]string{"regular", "admin", "suspended"})
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1), len(list2), "List of people not updated correctly")
}

func TestNewInfoJunkPerson(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("topsecret"), bcrypt.DefaultCost)
	require.Nil(t, err, "err should be nothing")

	list1, err := ListPeople(AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = NewPerson("John", "Buutler", "butler@ghj.co.uk", hashedPassword, false).Save("butler")
	require.Nil(t, err, "err should be nothing")

	list2, err := ListPeople([]string{"regular", "admin", "suspended"})
	require.Nil(t, err, "err should be nothing")

	require.Equal(t, len(list1)+1, len(list2), "List of people not updated correctly")
}

func TestSavePerson(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	list1, err := ListPeople(AllRoles)
	require.Nil(t, err, "err should be nothing")

	err = NewPerson("fred", "smith", junkEmail, junkHashedPassword, false).Save("fred")
	require.Nil(t, err, "err should be nothing")

	list2, err := ListPeople(AllRoles)
	require.Nil(t, err, "err should be nothing")
	assert.Equal(t, len(list1)+1, len(list2))
}

func TestAddPerson(t *testing.T) {
	r := require.New(t)

	err := ClearPeople()
	r.Nil(err, "err should be nothing")

	err = NewPerson("fred", "smith", junkEmail, junkHashedPassword, false).Save("fred")
	r.Nil(err, "err should be nothing")

	err = NewPerson("bloggs", "grey", junkEmail, junkHashedPassword, true).Save("bloggs")
	r.Nil(err, "err should be nothing")

	list, err := ListPeople(AllRoles)
	assert.Equal(t, 2, len(list))
	r.Nil(err, "err should be nothing")

	err = NewPerson("harry", "silver", junkEmail, junkHashedPassword, true).Save("harry")
	r.Nil(err, "err should be nothing")

	list, err = ListPeople([]string{"regular", "admin", "suspended"})
	assert.Equal(t, 3, len(list))
}

func TestPerson(t *testing.T) {
	r := require.New(t)

	// Clear all people
	err := ClearPeople()
	r.Nil(err, "err should be nothing")

	// Add a couple of people
	count := 0
	count = count + 1
	err = NewPerson("Fred", "xxx", junkEmail, junkHashedPassword, false).Save("1")
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = NewPerson("Bloggs", "xxx", junkEmail, junkHashedPassword, false).Save("2")
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = NewPerson("Jane", "xxx", junkEmail, junkHashedPassword, false).Save("3")
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = NewPerson("Alice", "xxx", junkEmail, junkHashedPassword, false).Save("4")
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = NewPerson("Bob", "xxx", junkEmail, junkHashedPassword, false).Save("5")
	r.Nil(err, "err should be nothing")

	// Check the expected number of People have been created
	list, err := ListPeople([]string{"regular", "admin", "suspended"})
	r.Nil(err, "err should be nothing")
	r.Equal(count, len(list), fmt.Sprintf("Unexpected number of people created. expected:%d actual:%d", count, len(list)))

	// Check the expected People have been created
	for _, id := range list {
		p, err := LoadPerson(id)
		r.Nil(err, "err should be nothing")

		found := false
		for _, id2 := range list {
			p2, err := LoadPerson(id2)
			r.Nil(err, "err should be nothing")

			equal := true
			if p.FirstName != p2.FirstName {
				equal = false
			}
			if p.LastName != p2.LastName {
				equal = false
			}
			if p.Email != p2.Email {
				equal = false
			}
			if p.Role != p2.Role {
				equal = false
			}
			if p.Player != p2.Player {
				equal = false
			}

			// Have we found the person?
			if equal {
				found = true
				break
			}
		}
		assert.Equal(t, found, true, fmt.Sprintf("Person [%s] not found", id))
	}

	// Delete the list of people
	for _, id := range list {
		err := RemovePerson(id)
		r.Nil(err, "err should be nothing")
	}

	// Check there are no more people
	list, err = ListPeople([]string{"regular", "admin", "suspended"})
	r.Nil(err, "err should be nothing")
	r.Equal(0, len(list), fmt.Sprintf("Unexpected number of people. Expected:%d, actual:%d", 0, len(list)))
}

func TestDeletePersonWithDuffID(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Clear the people
	err := ClearPeople()
	r.Nil(err, "err should be nothing")

	// Attempt to delete a person using a duff ID
	err = RemovePerson("junk")
	if err == nil {
		r.Fail(fmt.Sprintf("Expected an error. actually got: [%v].", err))
	} else {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusNotFound {
				r.Fail(fmt.Sprintf("Unexpected error: [%v]", err))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: [%v]", err))
		}
	}
}

func TestListPeopleWithDuffPlayerFile(t *testing.T) {
	r := require.New(t)

	// Clear the people
	err := ClearPeople()
	r.Nil(err, "err should be nothing")

	filename, err := makePersonFilename("x")
	r.Nil(err, "err should be nothing")

	// Create a new person file with junk contents
	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Attempt to use the 'junk' person!
	_, err = ListPeople([]string{"regular", "admin", "suspended"})
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				r.Fail(fmt.Sprintf("Unexpected error type: expected: %d, Actual: %d", http.StatusNotFound, cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("%s", err))
		}
	} else {
		r.Fail("Unexpected success")
	}
}

func TestDetailsWithDuffPersonFile(t *testing.T) {
	r := require.New(t)

	// Clear the people
	err := ClearPeople()
	r.Nil(err, "err should be nothing")

	// Create a new person file with junk contents
	filename, err := makePersonFilename("0")
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Check that List returns an error
	expected := "invalid character 'j' looking for beginning of value"
	_, err = LoadPerson("0")
	if err == nil {
		r.Fail(fmt.Sprintf("Error actual = (nil), and Expected = [%v].", expected))
	}
	if err.Error() != expected {
		r.Fail(fmt.Sprintf("Error actual = [%v], and Expected = [%v].", err, expected))
	}
}
