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
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestGetQueue(t *testing.T) {

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
			expectedStatus:        http.StatusOK,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "Bad token",
			token:                 "junk",
			expectedStatus:        http.StatusUnauthorized,
			expectedResultName:    "",
			expectedResultPlayers: []string{},
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetQueueRequest{
				Token: test.token,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/queue", bytes.NewBuffer(requestBody))
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
				var response GetQueueResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				actualName := response.Queue.Container.Name
				require.Equal(t, test.expectedResultName, actualName, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultName, actualName))

				actualPlayers := response.Queue.Container.Players
				if common.EqualArrayOfStrings(actualPlayers, test.expectedResultPlayers) {
					require.Fail(t, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultPlayers, actualPlayers))
				}
			}
		})
	}
}
