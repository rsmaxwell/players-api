package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestDeleteCourt(t *testing.T) {

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
		courtID        string
		expectedStatus int
	}{
		{
			testName:       "Good  request",
			token:          goodToken,
			courtID:        goodCourtID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			courtID:        goodCourtID,
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

			requestBody, err := json.Marshal(ListCourtsRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request
			req, err := http.NewRequest("DELETE", "/court/"+test.courtID, bytes.NewBuffer(requestBody))
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
				require.Equal(t, initialNumberOfCourts, finalNumberOfCourts+1, "Court was not deleted")
			} else {
				require.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
