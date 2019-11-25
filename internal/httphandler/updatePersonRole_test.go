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

func TestUpdatePersonRole(t *testing.T) {

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
		role           string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLoginCookie: true,
			sid:            goodSID,
			id:             anotherUserID,
			role:           person.RoleNormal,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "no login cookie",
			setLoginCookie: false,
			sid:            goodSID,
			id:             anotherUserID,
			role:           person.RoleNormal,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			setLoginCookie: true,
			sid:            goodSID,
			id:             "junk",
			role:           person.RoleNormal,
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "Bad Role",
			setLoginCookie: true,
			sid:            goodSID,
			id:             anotherUserID,
			role:           "junk",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "Admin Role",
			setLoginCookie: true,
			sid:            goodSID,
			id:             anotherUserID,
			role:           person.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Normal Role",
			setLoginCookie: true,
			sid:            goodSID,
			id:             anotherUserID,
			role:           person.RoleNormal,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Suspended Role",
			setLoginCookie: true,
			sid:            goodSID,
			id:             anotherUserID,
			role:           person.RoleSuspended,
			expectedStatus: http.StatusOK,
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
			requestBody, err := json.Marshal(UpdatePersonRoleRequest{
				Role: test.role,
			})
			require.Nil(t, err, "err should be nothing")

			req, err := http.NewRequest("PUT", contextPath+"/personrole/"+test.id, bytes.NewBuffer(requestBody))
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
				assert.Equal(t, p.Role, test.role, "The Person Role was not updated correctly")
			}
		})
	}
}
