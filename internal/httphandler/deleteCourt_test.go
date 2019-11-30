package httphandler

import (
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
	// * Login to get tokens
	// ***************************************************************
	accessTokenString, refreshTokenCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************

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
		courtID             string
		expectedStatus      int
	}{
		{
			testName:            "Good request",
			setAccessToken:      true,
			accessToken:         "Bearer " + accessTokenString,
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			courtID:             goodCourtID,
			expectedStatus:      http.StatusOK,
		},
		{
			testName:            "no login cookie",
			setAccessToken:      false,
			accessToken:         "",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			courtID:             goodCourtID,
			expectedStatus:      http.StatusUnauthorized,
		},
		{
			testName:            "bad Token",
			setAccessToken:      true,
			accessToken:         "junk",
			useGoodRefreshToken: true,
			setRefreshToken:     false,
			refreshToken:        "",
			courtID:             goodCourtID,
			expectedStatus:      http.StatusBadRequest,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfCourts, err := court.Size()
			require.Nil(t, err, "err should be nothing")

			// Set up the handlers on the router
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()

			// Create a request
			req, err := http.NewRequest("DELETE", contextPath+"/court/"+test.courtID, nil)
			require.Nil(t, err, "err should be nothing")

			setAccessToken(req, test.setAccessToken, test.accessToken)
			setRefreshToken(req, test.useGoodRefreshToken, test.setRefreshToken, refreshTokenCookie, test.refreshToken)

			// Serve the request
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
