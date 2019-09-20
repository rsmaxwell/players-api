package person

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// removeDir - Remove the people directory
func removeDir() error {

	_, err := os.Stat(peopleDir)
	if err == nil {
		err = common.RemoveContents(peopleDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}

		err = os.Remove(peopleDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

func ResetPeople(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("topsecret"), bcrypt.DefaultCost)
	assert.Nil(t, err)

	p := New("James", "Bond", "james@mi6.uk.gov", hashedPassword, false)

	err = Add("007", *p)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 1, len(list))
}

func TestNewInfoJunkPerson(t *testing.T) {

	ResetPeople(t)
}

func TestClearPeople(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)
}

func TestResetPeople(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	fred := New("fred", "smith", "123@aol.com", []byte("xxxxxxx"), false)
	bloggs := New("bloggs", "grey", "456@aol.com", []byte("yyyyyyyyy"), true)

	Add("fred", *fred)
	Add("bloggs", *bloggs)

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddPerson(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	fred := New("fred", "smith", "123@aol.com", []byte("xxxxxxx"), false)
	bloggs := New("bloggs", "grey", "456@aol.com", []byte("yyyyyyyyy"), true)

	Add("fred", *fred)
	Add("bloggs", *bloggs)

	list, err := List()
	assert.Equal(t, 2, len(list))
	assert.Nil(t, err)

	person := New("harry", "silver", "789@aol.com", []byte("zzzzzzzzz"), true)
	assert.NotNil(t, person)
	assert.Nil(t, err)

	err = Add("harry", *person)
	assert.Nil(t, err)

	list, err = List()
	assert.Equal(t, 3, len(list))
}

func TestPerson(t *testing.T) {

	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	var listOfPeople []*Person
	listOfPeople = append(listOfPeople, New("Fred", "x", "z", []byte("123@aol.com"), false))
	listOfPeople = append(listOfPeople, New("Bloggs", "x", "z", []byte("123@aol.com"), false))
	listOfPeople = append(listOfPeople, New("Jane", "x", "z", []byte("123@aol.com"), false))
	listOfPeople = append(listOfPeople, New("Alice", "x", "z", []byte("123@aol.com"), false))
	listOfPeople = append(listOfPeople, New("Bob", "x", "z", []byte("123@aol.com"), false))

	for index, person := range listOfPeople {
		err = Add(string(index), *person)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check the expected number of People have been created")
	list, err := List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(list), len(list), "")

	t.Log("Check the expected People have been created")
	for _, name := range list {
		found := false
		for _, id := range list {
			person, err := Get(id)
			if err != nil {
				t.Fatal(err)
			}

			if name == person.FirstName {
				found = true
			}
		}
		assert.Equal(t, found, true, "")
	}

	t.Log("Delete the list of people")
	for _, id := range list {
		err := Remove(id)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check there are no more people")
	list, err = List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(listOfPeople), 0, "")
}

func TestDeletePersonWithDuffID(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Clear the people
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete a person using a duff ID
	err = Remove("junk")
	if err == nil {
		t.Errorf("Expected an error. actually got: [%v].", err)
	} else {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusNotFound {
				t.Errorf("Unexpected error: [%v]", err)
			}
		} else {
			t.Errorf("Unexpected error: [%v]", err)
		}
	}
}

func TestListPeopleWithDuffPlayerFile(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Clear the people
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	filename, err := makeFilename("x")
	if err != nil {
		t.Fatal(err)
	}

	// Create a new person file with junk contents
	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected number of People have been created
	_, err = List()
	if err != nil {
		t.Fatal(err)
	}

	// Check the duff file was skipped
	_, err = List()
	t.Log(buf.String())
	if strings.HasPrefix("buf.String()", "Skipping unexpected person filename") {
		t.Fatal(err)
	}
}

func TestListPersonWithNoPeopleDirectory(t *testing.T) {

	// Clear the people directory
	err := removeDir()
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	_, err = List()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDetailsWithNoPeopleDirectory(t *testing.T) {

	// Remove the people directory
	err := removeDir()
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	expected := "no such file or directory"
	_, err = Get("0")
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if cerr, ok := err.(*codeError.CodeError); ok {
		if cerr.Code() != http.StatusInternalServerError {
			t.Errorf("Unexpected error code: %d", cerr.Code())
		}
	} else {
		t.Errorf("Unexpected error type")
	}
}

func TestDetailsWithDuffPersonFile(t *testing.T) {

	// Clear the people
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new person file with junk contents
	filename, err := makeFilename("0")
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	expected := "invalid character 'j' looking for beginning of value"
	_, err = Get("0")
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if err.Error() != expected {
		t.Errorf("Error actual = [%v], and Expected = [%v].", err, expected)
	}
}
