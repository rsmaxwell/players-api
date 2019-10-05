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

	"github.com/rsmaxwell/players-api/internal/model"
)

func TestUpdatePersonPlayer(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	// ***************************************************************
	// * Login to get a valid token
	// ***************************************************************
	goodToken, err := getLoginToken(t, goodUserID, goodPassword)
	require.Nil(t, err, "err should be nothing")

	// ***************************************************************
	// * Testcases
	// ***************************************************************
	tests := []struct {
		testName       string
		token          string
		id             string
		person         map[string]interface{}
		expectedStatus int
	}{
		{
			testName: "Good request",
			token:    goodToken,
			id:       goodUserID,
			person: map[string]interface{}{
				"FirstName": "aaa",
				"Lastname":  "bbb",
				"Email":     "123.456@xxx.com",
				"Player":    false,
				"Password":  "amother",
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Bad token",
			token:    "junk",
			id:       goodUserID,
			person: map[string]interface{}{
				"FirstName": "aaa",
				"Lastname":  "bbb",
				"Email":     "123.456@xxx.com",
				"Player":    false,
				"Password":  "amother",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// ***************************************************************
	// * Run the tests
	// ***************************************************************
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(UpdatePersonRequest{
				Token:  test.token,
				Person: test.person,
			})
			require.Nil(t, err, "err should be nothing")

			// Create a request to pass to our handler.
			req, err := http.NewRequest("PUT", "/person/"+test.id, bytes.NewBuffer(requestBody))
			require.Nil(t, err, "err should be nothing")

			// Pass the request to our handler
			router := mux.NewRouter()
			SetupHandlers(router)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, test.expectedStatus, rw.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus))

			// Check the person was actually updated
			if rw.Code == http.StatusOK {
				person, err := model.LoadPerson(test.id)
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
					assert.Equal(t, person.FirstName, value, "The Person lastname was not updated correctly")
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
