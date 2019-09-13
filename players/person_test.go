package players

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewRegisterRequest initialises a RegisterRequest object
func NewRegisterRequest(userID string, firstname string, lastname string, password string) (*RegisterRequest, error) {
	reg := new(RegisterRequest)
	reg.UserID = userID
	reg.FirstName = firstname
	reg.LastName = lastname
	reg.Password = password
	return reg, nil
}

func ResetPeople(t *testing.T) {

	err = ClearPeople()
	assert.Nil(t, err)

	err = RegisterPerson(NewRegisterRequest("007", "James", "Bond", "topsecret"))
	assert.Nil(t, err)

	err = RegisterPerson(NewPerson("fred", "Fred", "Bloggs", "qwerty"))
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddPerson(t *testing.T) {

	ResetPeople(t)

	err = RegisterPerson(NewRegisterRequest("007", "James", "Bond", "topsecret"))
	assert.Nil(t, err)

	count, err := CountPeople()
	assert.Nil(t, err)
	assert.Equal(t, 3, count)

	err = RegisterPerson(NewRegisterRequest("bob", "Robert", "Bruce", "haggis"))
	assert.Nil(t, err)

	count, err = CountPeople()
	assert.Nil(t, err)
	assert.Equal(t, 4, count)
}

func TestNewInfoJunkPerson(t *testing.T) {

	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	err = CreatePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	err = writefile(peopleInfoFile, "junk")
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		_, err = CreatePeopleInfoFile()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestNewInfoUnreadableInfofilePerson(t *testing.T) {

	// Remove all the contents of the person application directory
	t.Logf("Remove all the contents of the person application directory")
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new  "infofile"
	t.Logf("Create a new \"infofile\"")
	CreatePeopleInfoFile()

	t.Logf("Make the file \"%s\" unreadable", peopleInfoFile)
	err = os.Chmod(peopleInfoFile, 0000)
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		_, err := CreatePeopleInfoFile()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestPeople(t *testing.T) {

	t.Logf("Remove all the contents of the people directory")
	err := ClearPeople()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Create a number of new People")
	RegisterPerson(NewRegisterRequest("1", "aaa", "AAA", "one"))
	RegisterPerson(NewRegisterRequest("2", "bbb", "BBB", "two"))
	RegisterPerson(NewRegisterRequest("3", "ccc", "CCC", "three"))
	RegisterPerson(NewRegisterRequest("4", "ddd", "DDD", "four"))
	RegisterPerson(NewRegisterRequest("5", "eee", "EEE", "five"))

	t.Log("Check the expected number of People have been created")
	count, err := CountPeople()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, COUNT, "")
}
