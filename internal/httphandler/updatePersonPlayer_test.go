package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/model"
)

func TestUpdatePersonPlayer(t *testing.T) {

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
		testName       string
		setLoginCookie bool
		sid            string
		id             string
		player         bool
		expectedStatus int
	}{
		{
			testName:       "Good request - false",
			setLoginCookie: true,
			sid:            goodSID,
			id:             goodUserID,
			player:         false,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Good request - true",
			setLoginCookie: true,
			sid:            goodSID,
			id:             goodUserID,
			player:         true,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "no login cookie",
			setLoginCookie: false,
			sid:            goodSID,
			id:             goodUserID,
			player:         false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			setLoginCookie: true,
			sid:            goodSID,
			id:             "junk",
			player:         false,
			expectedStatus: http.StatusNotFound,
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
			requestBody, err := json.Marshal(UpdatePersonPlayerRequest{
				Player: test.player,
			})
			require.Nil(t, err, "err should be nothing")

			req, err := http.NewRequest("PUT", contextPath+"/personplayer/"+test.id, bytes.NewBuffer(requestBody))
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

			// Check the person was actually updated
			if rw.Code == http.StatusOK {
				p, err := person.Load(test.id)
				require.Nil(t, err, "err should be nothing")
				assert.Equal(t, p.Player, test.player, "The Person Player was not updated correctly")
			}
		})
	}
}
