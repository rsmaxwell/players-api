package httphandler

import (
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
	// * Login to get tokens
	// ***************************************************************
	accessTokenString, refreshTokenCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName              string
		setAccessToken        bool
		accessToken           string
		useGoodRefreshToken   bool
		setRefreshToken       bool
		refreshToken          string
		userID                string
		expectedStatus        int
		expectedResultName    string
		expectedResultPlayers []string
	}{
		{
			testName:              "Good request",
			setAccessToken:        true,
			accessToken:           "Bearer " + accessTokenString,
			useGoodRefreshToken:   true,
			setRefreshToken:       false,
			refreshToken:          "",
			expectedStatus:        http.StatusOK,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "no login cookie",
			setAccessToken:        false,
			accessToken:           "",
			useGoodRefreshToken:   true,
			setRefreshToken:       false,
			refreshToken:          "",
			expectedStatus:        http.StatusUnauthorized,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "bad token",
			setAccessToken:        true,
			accessToken:           "junk",
			useGoodRefreshToken:   true,
			setRefreshToken:       false,
			refreshToken:          "",
			expectedStatus:        http.StatusBadRequest,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()

			// Create a request
			req, err := http.NewRequest("GET", contextPath+"/queue", nil)
			require.Nil(t, err, "err should be nothing")

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
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
