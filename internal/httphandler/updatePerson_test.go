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

func TestUpdatePerson(t *testing.T) {

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
		person         map[string]interface{}
		expectedStatus int
	}{
		{
			testName:       "Good request",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             goodUserID,
			person: map[string]interface{}{
				"FirstName": "aaa",
				"LastName":  "bbb",
				"Email":     "123.456@xxx.com",
				"Player":    false,
				"Password":  "another",
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad userID",
			setLogonCookie: true,
			logonCookie:    logonCookie,
			id:             "junk",
			person: map[string]interface{}{
				"FirstName": "aaa",
				"LastName":  "bbb",
				"Email":     "123.456@xxx.com",
				"Player":    false,
				"Password":  "another",
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
			w := httptest.NewRecorder()

			// Create a request
			requestBody, err := json.Marshal(UpdatePersonRequest{
				Person: test.person,
			})
			require.Nil(t, err, "err should be nothing")

			r, err := http.NewRequest("PUT", contextPath+"/users/"+test.id, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			if test.setLogonCookie {
				r.AddCookie(test.logonCookie)
			}

			// Serve the request
			router.ServeHTTP(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", w.Code, test.expectedStatus))

			// Check the person was actually updated
			if w.Code == http.StatusOK {
				person, err := person.Load(test.id)
				require.Nil(t, err, "err should be nothing")

				if i, ok := test.person["FirstName"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.person[\"FirstName\"]' should be a string")
					}
					person.FirstName = value
					assert.Equal(t, person.FirstName, value, "The Person firstname was not updated correctly")
				}

				if i, ok := test.person["LastName"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.person[\"LastName\"]' should be a string")
					}
					person.LastName = value
					assert.Equal(t, person.LastName, value, "The Person lastname was not updated correctly")
				}

				if i, ok := test.person["Email"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.person[\"Email\"]' should be a string")
					}
					person.Email = value
					assert.Equal(t, person.Email, value, "The Person email was not updated correctly")
				}

				if i, ok := test.person["Player"]; ok {
					value, ok := i.(bool)
					if !ok {
						t.Errorf("The type of 'test.person[\"Player\"]' should be a boolean")
					}
					person.Player = value
					assert.Equal(t, person.Player, value, "The Person Player was not updated correctly")
				}
			}
		})
	}
}
