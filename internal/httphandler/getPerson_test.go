package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/model"
)

func TestGetPerson(t *testing.T) {

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
		testName              string
		token                 string
		userID                string
		expectedStatus        int
		expectedResultName    string
		expectedResultPlayers []string
	}{
		{
			testName:              "Good request",
			token:                 goodToken,
			userID:                goodUserID,
			expectedStatus:        http.StatusOK,
			expectedResultName:    "James",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "Bad token",
			token:                 "junk",
			userID:                goodUserID,
			expectedStatus:        http.StatusUnauthorized,
			expectedResultName:    "",
			expectedResultPlayers: []string{},
		},
		{
			testName:              "Bad userID",
			token:                 goodToken,
			userID:                "junk",
			expectedStatus:        http.StatusNotFound,
			expectedResultName:    "",
			expectedResultPlayers: []string{},
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetPersonRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", contextPath+"/person/"+test.userID, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			bytes, err := ioutil.ReadAll(rw.Body)
			require.Nil(t, err, "err should be nothing")

			if rw.Code == http.StatusOK {
				var response GetPersonResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, test.expectedResultName, response.Person.FirstName, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultName, response.Person.FirstName))
			}
		})
	}
}
