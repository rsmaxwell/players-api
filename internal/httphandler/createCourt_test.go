package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/peoplecontainer"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateCourt(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get a valid token
	// ***************************************************************
	goodToken, err := getLoginToken(t, goodUserID, goodPassword)
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		token          string
		name           string
		players        []string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          goodToken,
			name:           "Court 1",
			players:        []string{},
			expectedStatus: http.StatusOK,
		}, {
			testName:       "Bad token",
			token:          "junk",
			name:           "Court 1",
			players:        []string{},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfCourts, err := court.Size()
			require.Nil(t, err, "err should be nothing")

			requestBody, err := json.Marshal(CreateCourtRequest{
				Token: test.token,
				Court: court.Court{
					Container: peoplecontainer.PeopleContainer{
						Name:    test.name,
						Players: test.players,
					},
				},
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request
			req, err := http.NewRequest("POST", "/court", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			finalNumberOfCourts, err := court.Size()
			require.Nil(t, err, "err should be nothing")

			if rw.Code == http.StatusOK {
				require.Equal(t, initialNumberOfCourts+1, finalNumberOfCourts, "Court was not registered")
			} else {
				require.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
