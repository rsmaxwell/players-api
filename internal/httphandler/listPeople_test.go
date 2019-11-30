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
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestListPeople(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get tokens
	// ***************************************************************
	accessTokenString, refreshTokenCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Get a list of all the people
	// ***************************************************************
	allPeopleIDs, err := person.List(person.AllRoles)
	require.Nil(t, err, "err should be nothing")

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
		filter              []string
		expectedStatus      int
		expectedResult      []string
	}{
		{
			testName:            "Good request",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			filter:              person.AllRoles,
			expectedStatus:      http.StatusOK,
			expectedResult:      allPeopleIDs,
		},
		{
			testName:            "no login cookie",
			setAccessToken:      false,
			accessToken:         "",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			filter:              person.AllRoles,
			expectedStatus:      http.StatusUnauthorized,
			expectedResult:      allPeopleIDs,
		},
		{
			testName:            "bad token",
			setAccessToken:      true,
			accessToken:         "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			filter:              person.AllRoles,
			expectedStatus:      http.StatusBadRequest,
			expectedResult:      allPeopleIDs,
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
			requestBody, err := json.Marshal(ListPeopleRequest{
				Filter: test.filter,
			})
			require.Nil(t, err, "err should be nothing")

			req, err := http.NewRequest("GET", contextPath+"/users", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the response
			if rw.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(rw.Body)
				require.Nil(t, err, "err should be nothing")

				var response ListPeopleResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")

				// Check the response body is what we expect.
				if !common.EqualArrayOfStrings(response.People, test.expectedResult) {
					t.Logf("actual:   %v", response.People)
					t.Logf("expected: %v", test.expectedResult)
					t.Errorf("Unexpected list of people")
				}
			}
		})
	}
}
