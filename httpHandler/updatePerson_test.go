package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/person"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePerson(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		id             string
		person         map[string]interface{}
		expectedStatus int
	}{
		{
			testName: "Good request",
			token:    MyToken,
			id:       MyUserID,
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
			id:       MyUserID,
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

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(UpdatePersonRequest{
				Token:  test.token,
				Person: test.person,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("PUT", "/person/"+test.id, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can call its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, test.expectedStatus)
			}

			// Check the person was actually updated
			if rr.Code == http.StatusOK {
				person, err := person.Load(test.id)
				if err != nil {
					t.Fatal(err)
				}

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
