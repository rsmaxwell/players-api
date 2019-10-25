package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestDeletePerson(t *testing.T) {

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
		userID         string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          goodToken,
			userID:         anotherUserID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			userID:         anotherUserID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			token:          goodToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			require.Nil(t, err, "err should be nothing")

			requestBody, err := json.Marshal(ListPeopleRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request
			req, err := http.NewRequest("DELETE", contextPath+"/person/"+test.userID, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			finalNumberOfPeople, err := person.Size()
			require.Nil(t, err, "err should be nothing")

			if rw.Code == http.StatusOK {
				require.Equal(t, initialNumberOfPeople, finalNumberOfPeople+1, "Person was not deleted")
			} else {
				require.Equal(t, initialNumberOfPeople, finalNumberOfPeople, "Unexpected number of people")
			}
		})
	}
}
