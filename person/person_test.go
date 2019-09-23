package person

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// removeDir - Remove the people directory
func removeDir() error {

	_, err := os.Stat(baseDir)
	if err == nil {
		err = common.RemoveContents(baseDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}

		err = os.Remove(baseDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

func ResetPeople(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("topsecret"), bcrypt.DefaultCost)
	r.Nil(err, "err should be nothing")

	p := New("James", "Bond", "james@mi6.uk.gov", hashedPassword, false)

	err = Save("007", p)
	r.Nil(err, "err should be nothing")

	list, err := List()
	assert.Equal(t, 1, len(list))
}

func TestNewInfoJunkPerson(t *testing.T) {

	ResetPeople(t)
}

func TestClearPeople(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")
}

func TestResetPeople(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	err = Save("fred", New("fred", "smith", "123@aol.com", []byte("xxxxxxx"), false))
	r.Nil(err, "err should be nothing")

	err = Save("bloggs", New("bloggs", "grey", "456@aol.com", []byte("yyyyyyyyy"), true))
	r.Nil(err, "err should be nothing")

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddPerson(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	err = Save("fred", New("fred", "smith", "123@aol.com", []byte("xxxxxxx"), false))
	r.Nil(err, "err should be nothing")

	err = Save("bloggs", New("bloggs", "grey", "456@aol.com", []byte("yyyyyyyyy"), true))
	r.Nil(err, "err should be nothing")

	list, err := List()
	assert.Equal(t, 2, len(list))
	r.Nil(err, "err should be nothing")

	err = Save("harry", New("harry", "silver", "789@aol.com", []byte("zzzzzzzzz"), true))
	r.Nil(err, "err should be nothing")

	list, err = List()
	assert.Equal(t, 3, len(list))
}

func TestPerson(t *testing.T) {
	r := require.New(t)

	// Clear all people
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Add a couple of people
	count := 0
	count = count + 1
	err = Save("1", New("Fred", "x", "z", []byte("123@aol.com"), false))
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = Save("2", New("Bloggs", "x", "z", []byte("123@aol.com"), false))
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = Save("3", New("Jane", "x", "z", []byte("123@aol.com"), false))
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = Save("4", New("Alice", "x", "z", []byte("123@aol.com"), false))
	r.Nil(err, "err should be nothing")

	count = count + 1
	err = Save("5", New("Bob", "x", "z", []byte("123@aol.com"), false))
	r.Nil(err, "err should be nothing")

	// Check the expected number of People have been created
	list, err := List()
	r.Nil(err, "err should be nothing")
	r.Equal(count, len(list), fmt.Sprintf("Unexpected number of people created. expected:%d actual:%d", count, len(list)))

	// Check the expected People have been created
	for _, id := range list {
		p, err := Load(id)
		r.Nil(err, "err should be nothing")

		found := false
		for _, id2 := range list {
			p2, err := Load(id2)
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
			if p.Status != p2.Status {
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
		err := Remove(id)
		r.Nil(err, "err should be nothing")
	}

	// Check there are no more people
	list, err = List()
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
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Attempt to delete a person using a duff ID
	err = Remove("junk")
	if err == nil {
		r.Fail(fmt.Sprintf("Expected an error. actually got: [%v].", err))
	} else {
		if cerr, ok := err.(*codeError.CodeError); ok {
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
	err := Clear()
	r.Nil(err, "err should be nothing")

	filename, err := makeFilename("x")
	r.Nil(err, "err should be nothing")

	// Create a new person file with junk contents
	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Check the expected number of People have been created
	list, err := List()
	r.Nil(err, "err should be nothing")
	r.Equal(len(list), 1, "Unexpected number of people")

	// Attempt to use the 'junk' person!
	_, err = Load("junk")
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusNotFound {
				r.Fail(fmt.Sprintf("Unexpected error type: expected: %d, Actual: %d", http.StatusNotFound, cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("%s", err))
		}
	} else {
		r.Fail("Unexpected success")
	}
}

func TestListPersonWithNoPeopleDirectory(t *testing.T) {
	r := require.New(t)

	// Clear the people directory
	err := removeDir()
	r.Nil(err, "err should be nothing")

	// Attempt to use the people directory
	_, err = List()
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
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

func TestLoadWithNoPeopleDirectory(t *testing.T) {
	r := require.New(t)

	// Remove the people directory
	err := removeDir()
	r.Nil(err, "err should be nothing")

	// Attempt to use the people directory
	_, err = Load("0")
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusNotFound {
				r.Fail(fmt.Sprintf("Unexpected error type: expected: %d, Actual: %d", http.StatusNotFound, cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error type: expected: %s, Actual: %T", "*codeError.CodeError", err))
		}
	} else {
		r.Fail("Unexpected success")
	}
}

func TestDetailsWithDuffPersonFile(t *testing.T) {
	r := require.New(t)

	// Clear the people
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Create a new person file with junk contents
	filename, err := makeFilename("0")
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Check that List returns an error
	expected := "invalid character 'j' looking for beginning of value"
	_, err = Load("0")
	if err == nil {
		r.Fail(fmt.Sprintf("Error actual = (nil), and Expected = [%v].", expected))
	}
	if err.Error() != expected {
		r.Fail(fmt.Sprintf("Error actual = [%v], and Expected = [%v].", err, expected))
	}
}
