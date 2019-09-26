package court

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClearCourts(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	_, err = os.Stat(infoFile)
	r.Nil(err, "err should be nothing")
}

func TestResetCourt(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("fred"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("bloggs"))
	r.Nil(err, "err should be nothing")

	_, err = os.Stat(infoFile)
	r.Nil(err, "err should be nothing")

	list, err := List()
	r.Equal(2, len(list), fmt.Sprintf("Unexpected size of List: expected: %d, Actual: %d", 2, len(list)))
}

func TestAddCourt(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("fred"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("bloggs"))
	r.Nil(err, "err should be nothing")

	_, err = os.Stat(infoFile)
	r.Nil(err, "err should be nothing")

	list, err := List()
	r.Nil(err, "err should be nothing")
	r.Equal(2, len(list), fmt.Sprintf("Unexpected size of List: expected: %d, Actual: %d", 2, len(list)))

	_, err = Insert(New("harry"))
	r.Nil(err, "err should be nothing")

	list, err = List()
	r.Equal(3, len(list), fmt.Sprintf("Unexpected size of List: expected: %d, Actual: %d", 3, len(list)))
}

func TestNewInfoJunkCourt(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(infoFile, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")
}

func TestNewInfoUnreadableInfofileCourt(t *testing.T) {
	r := require.New(t)

	// Clear the court directory
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Make the court info file unreadable
	t.Logf("Make the file \"%s\" unreadable", infoFile)
	err = os.Chmod(infoFile, 0000)
	r.Nil(err, "err should be nothing")

	// Attempt to use the info file
	_, err = Insert(New("fred"))
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				r.Fail(fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: Expected = [*codeError.CodeError], Got = [%v}].", err))
		}
	} else {
		t.Errorf("Unexpected success")
	}
}

func TestGetAndIncrementCurrentIDCourt(t *testing.T) {
	r := require.New(t)

	t.Logf("Clear the court directory")
	err := Clear()
	r.Nil(err, "err should be nothing")

	for i := 0; i < 10; i++ {
		count, _ := getAndIncrementCurrentCourtID()
		assert.Equal(t, count, 1000+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofileCourt(t *testing.T) {
	r := require.New(t)

	t.Logf("Clear the court directory")
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Remove the court info file
	t.Logf("Remove the court info file")
	err = os.Remove(infoFile)
	r.Nil(err, "err should be nothing")

	assert.NotPanics(t, func() {
		getAndIncrementCurrentCourtID()
	})
}

func TestGetAndIncrementCurrentIDJunkContentsCourt(t *testing.T) {
	r := require.New(t)

	t.Logf("Clear the court directory")
	err := Clear()
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(infoFile, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	_, err = getAndIncrementCurrentCourtID()
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				r.Fail(fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: Expected = [*codeError.CodeError], Got = [%v}].", err))
		}
	} else {
		r.Fail("Unexpected success")
	}
}

func TestCourt(t *testing.T) {
	r := require.New(t)

	// Clear the court directory
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Create a number of new Courts
	_, err = Insert(New("Fred"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("Bloggs"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("Jane"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("Alice"))
	r.Nil(err, "err should be nothing")

	_, err = Insert(New("Bob"))
	r.Nil(err, "err should be nothing")

	// Check the expected number of Courts have been created
	list, err := List()
	r.Nil(err, "err should be nothing")
	assert.Equal(t, len(list), len(list), "")

	// Check the expected Courts have been created
	for _, id := range list {
		c, err := Load(id)
		r.Nil(err, "err should be nothing")

		found := false
		for _, id2 := range list {
			c2, err := Load(id2)
			r.Nil(err, "err should be nothing")

			equal := true
			if !container.Equal(c.Container, c2.Container) {
				equal = false
			}

			if equal {
				found = true
				break
			}
		}
		assert.Equal(t, found, true, fmt.Sprintf("Court [%s] not found", id))
	}

	// Delete the list of courts
	for _, id := range list {
		err := Remove(id)
		r.Nil(err, "err should be nothing")
	}

	// Check there are no more courts
	list, err = List()
	r.Nil(err, "err should be nothing")
	r.Equal(len(list), 0, "Unexpected number of courts")
}

func TestDeleteCourtWithDuffID(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Clear the courts
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Attempt to delete a court using a duff ID
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

func TestListCourtsWithDuffPlayerFile(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Clear the courts
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Create a new court file with junk contents
	err = ioutil.WriteFile(infoFile, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Attempt to use the court info file
	_, err = Insert(New("junk"))
	if err == nil {
		r.Fail(fmt.Sprintf("Expected an error. actually got: [%v].", err))
	} else {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				r.Fail(fmt.Sprintf("Unexpected error: [%v]", err))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: [%v]", err))
		}
	}
}

func TestLoadWithDuffCourtFile(t *testing.T) {
	r := require.New(t)

	// Clear the courts
	err := Clear()
	r.Nil(err, "err should be nothing")

	// Create a new court file with junk contents
	filename := listDir + "/junk.json"
	err = ioutil.WriteFile(filename, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	// Check that List returns an error
	_, err = Load("junk")
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				r.Fail(fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: Expected = [*codeError.CodeError], Got = [%v}].", err))
		}
	} else {
		r.Fail("Unexpected success")
	}
}
