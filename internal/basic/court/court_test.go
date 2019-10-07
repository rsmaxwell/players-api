package court

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SetupTestcase function
func setupTestcase(t *testing.T) func(t *testing.T) {

	list, err := List()
	require.Nil(t, err, "err should be nothing")

	for _, i := range list {
		err = Remove(i)
		require.Nil(t, err, "err should be nothing")
	}

	_, err = os.Stat(courtInfoFile)
	if err == nil {
		err = os.Remove(courtInfoFile)
		require.Nil(t, err, "err should be nothing")
	}

	return func(t *testing.T) {
	}
}

func TestNewInfoUnreadableInfofileCourt(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Make the court info file unreadable
	err := ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Attempt to use the info file
	_, err = New("Fred", []string{}).Add()
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				require.Fail(t, fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
			}
		} else {
			require.Fail(t, fmt.Sprintf("Unexpected error: Expected = [*codeerror.CodeError], Got = [%v}].", err))
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}

func TestGetAndIncrementCurrentIDCourt(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Count the initial number of courts
	list, err := List()
	require.Nil(t, err, "err should be nothing")

	for i := 0; i < 10; i++ {
		count, _ := getAndIncrementCurrentCourtID()
		require.Equal(t, count, 1000+len(list)+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofileCourt(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Remove the court info file
	_, err := os.Stat(courtInfoFile)
	if err == nil {
		err = os.Remove(courtInfoFile)
		require.Nil(t, err, "err should be nothing")
	}

	assert.NotPanics(t, func() {
		getAndIncrementCurrentCourtID()
	})
}

func TestGetAndIncrementCurrentIDJunkContentsCourt(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	err := ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	_, err = getAndIncrementCurrentCourtID()
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			require.Equal(t, cerr.Code(), http.StatusInternalServerError, fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
		} else {
			require.Fail(t, "Unexpected error: %s", err.Error())
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}

func TestListCourtsWithDuffPlayerFile(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Create a new court file with junk contents
	err := ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Attempt to use the court info file
	_, err = New("junk", []string{}).Add()
	if err == nil {
		require.Fail(t, fmt.Sprintf("Expected an error. actually got: [%v].", err))
	} else {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				require.Fail(t, fmt.Sprintf("Unexpected error: [%v]", err))
			}
		} else {
			require.Fail(t, fmt.Sprintf("Unexpected error: [%v]", err))
		}
	}
}

func TestLoadWithDuffCourtFile(t *testing.T) {
	teardown := setupTestcase(t)
	defer teardown(t)

	// Create a new court file with junk contents
	filename := courtListDir + "/junk.json"
	err := ioutil.WriteFile(filename, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Check that Load returns an error
	ref := common.Reference{Type: "court", ID: "junk"}
	_, err = Load(&ref)
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
			if cerr.Code() != http.StatusInternalServerError {
				require.Fail(t, fmt.Sprintf("Unexpected error code: %d", cerr.Code()))
			}
		} else {
			require.Fail(t, fmt.Sprintf("Unexpected error: Expected = [*codeerror.CodeError], Got = [%v}].", err))
		}
	} else {
		require.Fail(t, "Unexpected success")
	}
}
