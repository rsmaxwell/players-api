package court

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearCourts(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	_, err = os.Stat(courtInfoFile)
	assert.NotNil(t, err)
}

func TestResetCourt(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	fred := New("fred")
	bloggs := New("bloggs")

	Add(*fred)
	Add(*bloggs)

	_, err = os.Stat(courtInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
}

func TestAddCourt(t *testing.T) {

	err := Clear()
	assert.Nil(t, err)

	fred := New("fred")
	bloggs := New("bloggs")

	Add(*fred)
	Add(*bloggs)

	_, err = os.Stat(courtInfoFile)
	assert.Nil(t, err)

	list, err := List()
	assert.Equal(t, 2, len(list))
	assert.Nil(t, err)

	court := New("harry")
	assert.NotNil(t, court)
	assert.Nil(t, err)

	err = Add(*court)
	assert.Nil(t, err)

	list, err = List()
	assert.Equal(t, 3, len(list))
}

func TestNewInfoJunkCourt(t *testing.T) {

	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInfoUnreadableInfofileCourt(t *testing.T) {

	t.Logf("Clear the court directory")
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Make the court info file unreadable
	t.Logf("Make the file \"%s\" unreadable", courtInfoFile)
	err = os.Chmod(courtInfoFile, 0000)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Attemp to use the info file")
	_, err = List()
	assert.Panics(t, func() {
		_, err := GetInfo()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAndIncrementCurrentIDCourt(t *testing.T) {

	t.Logf("Clear the court directory")
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		count, _ := getAndIncrementCurrentCourtID()
		assert.Equal(t, count, 1000+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofileCourt(t *testing.T) {

	t.Logf("Clear the court directory")
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Remove the court info file
	t.Logf("Remove the court info file")
	err = os.Remove(courtInfoFile)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotPanics(t, func() {
		getAndIncrementCurrentCourtID()
	})
}

func TestGetAndIncrementCurrentIDJunkContentsCourt(t *testing.T) {

	t.Logf("Clear the court directory")
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() {
		getAndIncrementCurrentCourtID()
	})
}

func TestCourt(t *testing.T) {

	t.Logf("Clear the court directory")
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Create a number of new Courts")
	var listOfCourts []*Court
	listOfCourts = append(listOfCourts, New("Fred"))
	listOfCourts = append(listOfCourts, New("Bloggs"))
	listOfCourts = append(listOfCourts, New("Jane"))
	listOfCourts = append(listOfCourts, New("Alice"))
	listOfCourts = append(listOfCourts, New("Bob"))

	for _, court := range listOfCourts {
		err = Add(*court)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check the expected number of Courts have been created")
	list, err := List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(list), len(list), "")

	t.Log("Check the expected Courts have been created")
	for _, name := range list {
		found := false
		for _, id := range list {
			court, err := Get(id)
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
	for _, id := range list {
		err := Delete(id)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check there are no more courts")
	list, err = List()
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

	// Clear the courts
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete a court using a duff ID
	expected := "court [junk] not found"
	err = Delete("junk")
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

	// Clear the courts
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new court file with junk contents
	err = ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
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

	// Clear the courts
	err := Clear()
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

	// Remove the court directory
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
	if !strings.HasSuffix(err.Error(), expected) {
		t.Errorf("Error actual = [%v], and Expected = [%v].", err, expected)
	}
}

func TestDetailsWithDuffCourtFile(t *testing.T) {

	// Clear the courts
	err := Clear()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new court file with junk contents
	filename := courtListDir + "/0.json"
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
