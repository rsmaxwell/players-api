package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestGetQueue(t *testing.T) {

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
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName              string
		setLoginCookie        bool
		sid                   string
		userID                string
		expectedStatus        int
		expectedResultName    string
		expectedResultPlayers []string
	}{
		{
			testName:              "Good request",
			setLoginCookie:        true,
			sid:                   goodSID,
			expectedStatus:        http.StatusOK,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "no login cookie",
			setLoginCookie:        false,
			sid:                   goodSID,
			expectedStatus:        http.StatusUnauthorized,
			expectedResultName:    "Queue",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "bad sid",
			setLoginCookie:        true,
			sid:                   "junk",
			expectedStatus:        http.StatusUnauthorized,
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
