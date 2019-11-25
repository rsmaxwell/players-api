package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	// * Login to get valid session
	// ***************************************************************
	req, err := http.NewRequest("GET", contextPath+"/login", nil)
	require.Nil(t, err, "err should be nothing")

	userID := "007"
	password := "topsecret"
	req.Header.Set("Authorization", model.BasicAuth(userID, password))

	router := mux.NewRouter()
	SetupHandlers(router)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	sess, err := globalSessions.SessionStart(rw, req)
	require.Nil(t, err, "err should be nothing")
	defer sess.SessionRelease(rw)

	goodSID := sess.SessionID()
	require.NotNil(t, goodSID, "err should be nothing")

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
		setLoginCookie bool
		sid            string
		filter         []string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			setLoginCookie: true,
			sid:            goodSID,
			filter:         person.AllRoles,
			expectedStatus: http.StatusOK,
			expectedResult: allPeopleIDs,
		},
		{
			testName:       "no login cookie",
			setLoginCookie: false,
			sid:            goodSID,
			filter:         person.AllRoles,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: allPeopleIDs,
		},
		{
			testName:       "bad sid",
			setLoginCookie: true,
			sid:            "junk",
			filter:         person.AllRoles,
			expectedStatus: http.StatusUnauthorized,
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
			rw := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(ListPeopleRequest{
				Filter: test.filter,
			})
			require.Nil(t, err, "err should be nothing")

			req, err := http.NewRequest("GET", contextPath+"/person", bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// set a cookie with the value of the login sid
			if test.setLoginCookie {
				cookieLifeTime := 3 * 60 * 60
				cookie := http.Cookie{
					Name:    "players-api",
					Value:   test.sid,
					MaxAge:  cookieLifeTime,
					Expires: time.Now().Add(time.Duration(cookieLifeTime) * time.Second),
				}
				req.AddCookie(&cookie)
			}

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
