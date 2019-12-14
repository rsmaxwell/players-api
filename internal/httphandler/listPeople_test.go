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
	// * Login
	// ***************************************************************
	logonCookie := testLogin(t, "007", "topsecret")

	// ***************************************************************
	// * Get a list of all the people
	// ***************************************************************
	allPeopleIDs, err := person.List(person.AllRoles)
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		setLogonCookie bool
		logonCookie    *http.Cookie
		filter         []string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			filter:         person.AllRoles,
			expectedStatus: http.StatusOK,
			expectedResult: allPeopleIDs,
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
			w := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(ListPeopleRequest{
				Filter: test.filter,
			})
			require.Nil(t, err, "err should be nothing")

			r, err := http.NewRequest("POST", contextPath+"/users", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the response
			if w.Code == http.StatusOK {
				bytes, err := ioutil.ReadAll(w.Body)
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
