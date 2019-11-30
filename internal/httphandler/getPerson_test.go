package httphandler

import (
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
	// * Login to get tokens
	// ***************************************************************
	accessTokenString, refreshTokenCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName            string
		setAccessToken      bool
		accessToken         string
		useGoodRefreshToken bool
		setRefreshToken     bool
		refreshToken        string
		userID              string
		expectedStatus      int
		expectedResultName  string
	}{
		{
			testName:            "Good request",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              goodUserID,
			expectedStatus:      http.StatusOK,
			expectedResultName:  "James",
		},
		{
			testName:            "no login cookie",
			setAccessToken:      false,
			accessToken:         "",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              "junk",
			expectedStatus:      http.StatusUnauthorized,
			expectedResultName:  "",
		},
		{
			testName:            "bad token",
			setAccessToken:      true,
			accessToken:         "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              "junk",
			expectedStatus:      http.StatusBadRequest,
			expectedResultName:  "",
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Set up the handlers on the router
			router := mux.NewRouter()
			router2 := Middleware(router)
			SetupHandlers(router)
			rw := httptest.NewRecorder()

			// Create a request
			req, err := http.NewRequest("GET", contextPath+"/users/"+test.userID, nil)
			require.Nil(t, err, "err should be nothing")

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
			router2.ServeHTTP(rw, req)
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
