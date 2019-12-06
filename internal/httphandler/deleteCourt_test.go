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
	// * Login
	// ***************************************************************
	logonCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		courtID        string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			courtID:        goodCourtID,
			expectedStatus: http.StatusOK,
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
			w := httptest.NewRecorder()

			// Create a request
			r, err := http.NewRequest("DELETE", contextPath+"/court/"+test.courtID, nil)
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			finalNumberOfCourts, err := court.Size()
			require.Nil(t, err, "err should be nothing")

			if w.Code == http.StatusOK {
				require.Equal(t, initialNumberOfCourts, finalNumberOfCourts+1, "Court was not deleted")
			} else {
				require.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
