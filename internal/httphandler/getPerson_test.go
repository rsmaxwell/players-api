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
	// * Login
	// ***************************************************************
	logonCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName           string
		setLogonCookie bool
		logonCookie    *http.Cookie
		userID             string
		expectedStatus     int
		expectedResultName string
	}{
		{
			testName:           "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			userID:             goodUserID,
			expectedStatus:     http.StatusOK,
			expectedResultName: "James",
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
			w := httptest.NewRecorder()

			// Create a request
			r, err := http.NewRequest("GET", contextPath+"/users/"+test.userID, nil)
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router2.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			bytes, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err, "err should be nothing")

			if w.Code == http.StatusOK {
				var response GetPersonResponse
				err = json.Unmarshal(bytes, &response)
				require.Nil(t, err, "err should be nothing")
				require.Equal(t, test.expectedResultName, response.Person.FirstName, fmt.Sprintf("handler returned unexpected body: want %v, got %v", test.expectedResultName, response.Person.FirstName))
			}
		})
	}
}
