package model

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"

	"github.com/stretchr/testify/require"
)

func TestCourt(t *testing.T) {
	teardown := SetupFull(t)
	defer teardown(t)

	// Count the initial number of courts
	list1, err := court.List()
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
		id, err := court.New(i.name, i.players).Add()
		datacourts[index].id = id
		require.Nil(t, err, "err should be nothing")
	}

	// Check the expected number of Courts have been created
	list2, err := court.List()
	require.Nil(t, err, "err should be nothing")
	require.Equal(t, len(list1)+len(datacourts), len(list2), "Unexpected number of courts")

	// Check the Courts have been created correctly
	for _, i := range datacourts {
		ref := common.Reference{Type: "court", ID: i.id}
		c, err := court.Load(&ref)
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
	err := court.Remove("junk")
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

			ref := &common.Reference{Type: "court", ID: test.courtID}
			err := court.Update(ref, test.court)
			require.Nil(t, err, "err should be nothing")

			// Check the court was actually updated
			c, err := court.Load(ref)
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
