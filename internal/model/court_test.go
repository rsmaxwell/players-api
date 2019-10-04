package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInfoUnreadableInfofileCourt(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Make the court info file unreadable
	err := ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Attempt to use the info file
	_, err = NewCourt("Fred", []string{}).Add()
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
	teardown := SetupFull(t)
	defer teardown(t)

	// Count the initial number of courts
	list, err := ListCourts()
	require.Nil(t, err, "err should be nothing")

	for i := 0; i < 10; i++ {
		count, _ := getAndIncrementCurrentCourtID()
		require.Equal(t, count, 1000+len(list)+i, "Unexpected value of ID")
	}
}

func TestGetAndIncrementCurrentIDNoInfofileCourt(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Remove the court info file
	t.Logf("Remove the court info file")
	err := os.Remove(courtInfoFile)
	require.Nil(t, err, "err should be nothing")

	assert.NotPanics(t, func() {
		getAndIncrementCurrentCourtID()
	})
}

func TestGetAndIncrementCurrentIDJunkContentsCourt(t *testing.T) {
	teardown := SetupFull(t)
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

func TestCourt(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Count the initial number of courts
	list1, err := ListCourts()
	require.Nil(t, err, "err should be nothing")

	// Create a number of new Courts
	datacourts := []struct {
		id      string
		name    string
		players []string
	}{
		{name: "North", players: []string{}},
		{name: "West", players: []string{}},
		{name: "East", players: []string{}},
		{name: "West", players: []string{}},
		{name: "Center", players: []string{}},
	}

	for index, i := range datacourts {
		id, err := NewCourt(i.name, i.players).Add()
		datacourts[index].id = id
		require.Nil(t, err, "err should be nothing")
	}

	// Check the expected number of Courts have been created
	list2, err := ListCourts()
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1)+len(datacourts), len(list2), "Unexpected number of courts")

	// Check the Courts have been created correctly
	for _, i := range datacourts {
		ref := common.Reference{Type: "court", ID: i.id}
		c, err := LoadCourt(&ref)
		require.Nil(t, err, "err should be nothing")

		require.Equal(t, c.Container.Name, i.name, fmt.Sprintf("Court [%s] has the wrong name", i.id))

		if !common.EqualArrayOfStrings(c.Container.Players, i.players) {
			require.Fail(t, fmt.Sprintf("Court [%s] has the wrong players", i.id))
		}
	}
}

func TestDeleteCourtWithDuffID(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Attempt to delete a court using a duff ID
	err := RemoveCourt("junk")
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

func TestListCourtsWithDuffPlayerFile(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Create a new court file with junk contents
	err := ioutil.WriteFile(courtInfoFile, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Attempt to use the court info file
	_, err = NewCourt("junk", []string{}).Add()
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
	teardown := SetupFull(t)
	defer teardown(t)

	// Create a new court file with junk contents
	filename := courtListDir + "/junk.json"
	err := ioutil.WriteFile(filename, []byte("junk"), 0644)
	require.Nil(t, err, "err should be nothing")

	// Check that Load returns an error
	ref := common.Reference{Type: "court", ID: "junk"}
	_, err = LoadCourt(&ref)
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

func TestUpdateCourt(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		personID       string
		courtID        string
		court          map[string]interface{}
		expectedStatus int
	}{
		{
			testName: "Good request",
			personID: "007",
			courtID:  "1001",
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"bob", "jill", "alice"},
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			token, err := session.New(test.personID)
			require.Nil(t, err)

			mySession := session.LookupToken(token)
			require.Nil(t, err)

			ref := &common.Reference{Type: "court", ID: test.courtID}
			err = UpdateCourt(ref, mySession, test.court)
			require.Nil(t, err, "err should be nothing")

			// Check the court was actually updated
			c, err := LoadCourt(ref)
			if err != nil {
				require.Fail(t, err.Error())
			}

			if j, ok := test.court["Container"]; ok {

				container2, ok := j.(map[string]interface{})
				if !ok {
					require.Fail(t, fmt.Sprintf("The type of 'test.court[\"Name\"]' should be a 'map[string]interface{}'"))
				}

				if i, ok := container2["Name"]; ok {
					value, ok := i.(string)
					if !ok {
						require.Fail(t, fmt.Sprintf("The type of 'test.court[\"Container\"][\"Name\"]' should be a string"))
					}
					require.Equal(t, c.Container.Name, value, "The Court name was not updated correctly")
				}

				if i, ok := container2["Players"]; ok {
					value, ok := i.([]string)
					if !ok {
						require.Fail(t, fmt.Sprintf("The type of 'test.court[\"Container\"][\"Players\"]' should be an array of strings"))
					}
					if !common.EqualArrayOfStrings(c.Container.Players, value) {
						require.Fail(t, fmt.Sprintf("The Court name was not updated correctly:\n got %v\n want %v", c.Container.Players, value))
					}
				}
			}
		})
	}
}
