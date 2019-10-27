package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCourt(t *testing.T) {

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
		court          map[string]interface{}
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLoginCookie: true,
			sid:            goodSID,
			id:             goodCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"bob", "jill", "alice"},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "no login cookie",
			setLoginCookie: false,
			sid:            goodSID,
			id:             goodCourtID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "bad sid",
			setLoginCookie: true,
			sid:            "junk",
			id:             goodCourtID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			setLoginCookie: true,
			sid:            goodSID,
			id:             "junk",
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "Bad player",
			setLoginCookie: true,
			sid:            goodSID,
			id:             goodCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"junk"},
				},
			},
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
			requestBody, err := json.Marshal(UpdateCourtRequest{
				Court: test.court,
			})
			require.Nil(t, err, "err should be nothing")

			req, err := http.NewRequest("PUT", contextPath+"/court/"+test.id, bytes.NewBuffer(requestBody))
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
				ref := common.Reference{Type: "court", ID: test.id}
				c, err := court.Load(&ref)
				require.Nil(t, err, "err should be nothing")

				if i, ok := test.court["Name"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Name\"]' should be a string")
					}
					c.Container.Name = value
					assert.Equal(t, c.Container.Name, value, "The Court name was not updated correctly")
				}

				if i, ok := test.court["Players"]; ok {
					value, ok := i.([]string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Players\"]' should be an array of strings")
					}
					c.Container.Players = value
					if !common.EqualArrayOfStrings(c.Container.Players, value) {
						t.Errorf("The Court name was not updated correctly:\n got %v\n want %v", c.Container.Players, value)
					}
				}
			}
		})
	}
}
