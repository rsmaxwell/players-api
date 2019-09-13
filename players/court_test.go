package players

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveCourtDirectory(t *testing.T) {

	err := RemoveCourtDirectory()
	assert.Nil(t, err)

	_, err = os.Stat(courtInfoFile)
	assert.NotNil(t, err)
}

func TestResetCourt(t *testing.T) {

	fred = NewCourt("fred")
	bloggs = NewCourt("bloggs")

	err := Reset(fred, bloggs)
	assert.Nil(t, err)

	_, err = os.Stat(courtInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddCourt(t *testing.T) {

	fred = NewCourt("fred")
	bloggs = NewCourt("bloggs")

	err := Reset(fred, bloggs)
	assert.Nil(t, err)

	_, err = os.Stat(courtInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
	assert.Nil(t, err)

	court, err := NewCourt("harry")
	assert.NotNil(t, court)
	assert.Nil(t, err)

	err = AddCourt(*court)
	assert.Nil(t, err)

	list, err = List()
	assert.Equal(t, 3, len(list))
}

func TestNewInfoJunkCourt(t *testing.T) {

	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	err = CreateCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	err = writefile(courtInfoFile, "junk")
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		_, err = CreateCourtInfoFile()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestNewInfoUnreadableInfofileCourt(t *testing.T) {

	// Remove all the contents of the court application directory
	t.Logf("Remove all the contents of the court application directory")
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new  "infofile"
	t.Logf("Create a new \"infofile\"")
	CreateCourtInfoFile()

	t.Logf("Make the file \"%s\" unreadable", courtInfoFile)
	err = os.Chmod(courtInfoFile, 0000)
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		_, err := CreateCourtInfoFile()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAndIncrementCurrentIDCourt(t *testing.T) {

	// Remove all the contents of the court application directory
	t.Logf("Remove all the contents of the court application directory")
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new "infofile"
	t.Logf("Create a new \"infofile\"")
	_, err = CreateCourtInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		id, _ := GetAndIncrementCurrentID()
		assert.Equal(t, id, 1000+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofileCourt(t *testing.T) {

	// Remove all the contents of the court application directory
	t.Logf("Remove all the contents of the court application directory")
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotPanics(t, func() {
		GetAndIncrementCurrentID()
	})
}

func TestGetAndIncrementCurrentIDJunkContentsCourt(t *testing.T) {

	t.Logf("Remove all the contents of the court application directory")
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create the court directory")
	err = CreateCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create a court InfoFile with junk contents")
	err = writefile(courtInfoFile, "junk")
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		GetAndIncrementCurrentID()
	})
}

func TestCourt(t *testing.T) {

	t.Logf("Remove all the contents of the court directory")
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create a new \"infofile\"")
	_, err = CreateCourtInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Create a number of new Courts")
	listOfCourts := [...]Court{}
	listOfCourts = append(listOfCourts, NewCourt("Fred"))
	listOfCourts = append(listOfCourts, NewCourt("Bloggs"))
	listOfCourts = append(listOfCourts, NewCourt("Jane"))
	listOfCourts = append(listOfCourts, NewCourt("Alice"))
	listOfCourts = append(listOfCourts, NewCourt("Bob"))

	for i, court := range listOfNames {
		err = AddCourt(*court)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check the expected number of Courts have been created")
	listOfCourts, err := List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(listOfCourts), len(listOfNames), "")

	t.Log("Check the expected Courts have been created")
	for _, name := range listOfNames {
		found := false
		for _, id := range listOfCourts {
			court, err := Details(id)
			if err != nil {
				t.Fatal(err)
			}

			if name == court.Name {
				found = true
			}
		}
		assert.Equal(t, found, true, "")
	}

	t.Log("Delete the list of courts")
	for _, id := range listOfCourts {
		err := Delete(id)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check there are no more courts")
	listOfCourts, err = List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(listOfCourts), 0, "")
}

func TestDeleteCourtWithDuffID(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Remove all the contents of the courts application directory
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete a court using a duff ID
	expected := "court [9999999] not found"
	err = Delete(9999999)
	if err == nil {
		t.Errorf("Error actual = (nil), and Expected = [%v].", expected)
	}
	if err.Error() != expected {
		t.Errorf("Error actual = [%v], and Expected = [%v.]", err, expected)
	}
}

func TestListCourtsWithDuffPlayerFile(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Remove all the contents of the courts application directory
	err := RemoveCourtDirectory()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new infofile
	_, err = CreateCourtInfoFile()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new court file with junk contents
	err = writefileInDirectory(courtDataDirectory, "not-a-number", "junk")
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected number of Courts have been created
	_, err = List()
	if err != nil {
		t.Fatal(err)
	}

	// Check the duff file was skipped
	_, err = List()
	t.Log(buf.String())
	if strings.HasPrefix("buf.String()", "Skipping unexpected court filename") {
		t.Fatal(err)
	}
}

func TestListCourtsWithNoCourtsDirectory(t *testing.T) {

	// Remove the contents of the courts directory
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

func TestDetailsWithNoCourtDirectory(t *testing.T) {

	// Remove the contents of the court directory
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

func TestDetailsWithDuffCourtFile(t *testing.T) {

	// Remove the court directory
	err := Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new court file with junk contents
	err = writefileInDirectory(courtDataDirectory, "0.json", "junk")
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
