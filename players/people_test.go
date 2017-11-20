package players

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writefile(filepath string, contents string) error {

	data := []byte(contents)
	err := ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func writefileInDirectory(directory string, filename string, contents string) error {
	return writefile(directory+"/"+filename, contents)
}

func TestRemovePersonDirectory(t *testing.T) {

	err := RemovePeopleDirectory()
	assert.Nil(t, err)

	_, err = os.Stat(peopleInfoFile)
	assert.NotNil(t, err)
}

func TestReset(t *testing.T) {

	err := Reset("fred", "bloggs")
	assert.Nil(t, err)

	_, err = os.Stat(peopleInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddPerson(t *testing.T) {

	err := Reset("fred", "bloggs")
	assert.Nil(t, err)

	_, err = os.Stat(peopleInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
	assert.Nil(t, err)

	person, err := NewPerson("harry")
	assert.NotNil(t, person)
	assert.Nil(t, err)

	err = AddPerson(*person)
	assert.Nil(t, err)

	list, err = List()
	assert.Equal(t, 3, len(list))
}

func TestNewInfoJunk(t *testing.T) {

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

func TestNewInfoUnreadableInfofile(t *testing.T) {

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

func TestGetAndIncrementCurrentID(t *testing.T) {

	// Remove all the contents of the people application directory
	t.Logf("Remove all the contents of the people application directory")
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new "infofile"
	t.Logf("Create a new \"infofile\"")
	_, err = CreatePeopleInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		id, _ := GetAndIncrementCurrentID()
		assert.Equal(t, id, 1000+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofile(t *testing.T) {

	// Remove all the contents of the person application directory
	t.Logf("Remove all the contents of the person application directory")
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotPanics(t, func() {
		GetAndIncrementCurrentID()
	})
}

func TestGetAndIncrementCurrentIDJunkContents(t *testing.T) {

	t.Logf("Remove all the contents of the people application directory")
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create the person directory")
	err = CreatePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create a person InfoFile with junk contents")
	err = writefile(peopleInfoFile, "junk")
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		GetAndIncrementCurrentID()
	})
}

func TestPeople(t *testing.T) {

	t.Logf("Remove all the contents of the people directory")
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create a new \"infofile\"")
	_, err = CreatePeopleInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Create a number of new People")
	listOfNames := [...]string{"Fred", "Bloggs", "Jane", "Alice", "Bob"}

	for i, name := range listOfNames {
		// Create a new  "infofile"
		t.Logf("(%d) Create a new Player [%s]", i, name)

		p, err := NewPerson(name)
		if err != nil {
			t.Fatal(err)
		}

		err = AddPerson(*p)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check the expected number of People have been created")
	listOfPeople, err := List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(listOfPeople), len(listOfNames), "")

	t.Log("Check the expected People have been created")
	for _, name := range listOfNames {
		found := false
		for _, id := range listOfPeople {
			person, err := Details(id)
			if err != nil {
				t.Fatal(err)
			}

			if name == person.Name {
				found = true
			}
		}
		assert.Equal(t, found, true, "")
	}

	t.Log("Delete the list of people")
	for _, id := range listOfPeople {
		err := Delete(id)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check there are no more people")
	listOfPeople, err = List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(listOfPeople), 0, "")
}

func TestDeletePlayerWithDuffID(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Remove all the contents of the people application directory
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete a person using a duff ID
	expected := "person [9999999] not found"
	err = Delete(9999999)
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if err.Error() != expected {
		t.Errorf("Error actual = [%v], and Expected = [%v.]", err, expected)
	}
}

func TestListPlayersWithDuffPlayerFile(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Remove all the contents of the people application directory
	err := RemovePeopleDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new infofile
	_, err = CreatePeopleInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new person file with junk contents
	err = writefileInDirectory(peopleDataDirectory, "not-a-number", "junk")
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected number of Players have been created
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

func TestListPlayersWithNoPlayerDirectory(t *testing.T) {

	// Remove the contents of the people directory
	err := Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	_, err = List()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDetailsWithNoPlayerDirectory(t *testing.T) {

	// Remove the contents of the people directory
	err := Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	expected := "no such file or directory"
	_, err = Details(0)
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if !strings.HasSuffix(err.Error(), expected) {
		t.Errorf("Error actual = [%v], and Expected = [%v].", err, expected)
	}
}

func TestDetailsWithDuffPlayerFile(t *testing.T) {

	// Remove the people directory
	err := Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new person file with junk contents
	err = writefileInDirectory(peopleDataDirectory, "0.json", "junk")
	if err != nil {
		t.Fatal(err)
	}

	// Check that List returns an error
	expected := "invalid character 'j' looking for beginning of value"
	_, err = Details(0)
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if err.Error() != expected {
		t.Errorf("Error actual = [%v], and Expected = [%v].", err, expected)
	}
}
