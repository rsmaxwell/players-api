package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
		id             string
		role           string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             anotherUserID,
			role:           person.RoleNormal,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad userID",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             "junk",
			role:           person.RoleNormal,
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "Bad Role",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             anotherUserID,
			role:           "junk",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "Admin Role",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             anotherUserID,
			role:           person.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Normal Role",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             anotherUserID,
			role:           person.RoleNormal,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Suspended Role",
			setLogonCookie: true,
			logonCookie:    logonCookie,
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
			w := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(UpdatePersonRoleRequest{
				Role: test.role,
			})
			require.Nil(t, err, "err should be nothing")

			r, err := http.NewRequest("PUT", contextPath+"/users/role/"+test.id, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the person was actually updated
			if w.Code == http.StatusOK {
				p, err := person.Load(test.id)
				require.Nil(t, err, "err should be nothing")
				assert.Equal(t, p.Role, test.role, "The Person Role was not updated correctly")
			}
		})
	}
}
