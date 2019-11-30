package httphandler

import (
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
	}{
		{
			testName:            "Good request",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              anotherUserID,
			expectedStatus:      http.StatusOK,
		},
		{
			testName:            "no login cookie",
			setAccessToken:      false,
			accessToken:         "",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              anotherUserID,
			expectedStatus:      http.StatusUnauthorized,
		},
		{
			testName:            "bad token",
			setAccessToken:      true,
			accessToken:         "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              anotherUserID,
			expectedStatus:      http.StatusBadRequest,
		},
		{
			testName:            "Bad userID",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              "junk",
			expectedStatus:      http.StatusNotFound,
		},
		{
			testName:            "delete myself",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			userID:              "007",
			expectedStatus:      http.StatusUnauthorized,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			require.Nil(t, err, "err should be nothing")

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()

			// Create a request
			req, err := http.NewRequest("DELETE", contextPath+"/users/"+test.userID, nil)
			require.Nil(t, err, "err should be nothing")

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
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
