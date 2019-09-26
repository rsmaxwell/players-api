package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/utilities"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCourt(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		id             string
		court          map[string]interface{}
		expectedStatus int
	}{
		{
			testName: "Good request",
			token:    MyToken,
			id:       MyCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"bob", "jill", "alice"},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Bad token",
			token:    "junk",
			id:       MyCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{},
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName: "Bad userID",
			token:    MyToken,
			id:       "junk",
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			testName: "Bad player",
			token:    MyToken,
			id:       MyCourtID,
			court: map[string]interface{}{
				"Container": map[string]interface{}{
					"Name":    "COURT 101",
					"Players": []string{"junk"},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(UpdateCourtRequest{
				Token: test.token,
				Court: test.court,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("PUT", "/court/"+test.id, bytes.NewBuffer(requestBody))
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

			// Check the court was actually updated
			if rr.Code == http.StatusOK {
				court, err := court.Load(test.id)
				if err != nil {
					t.Fatal(err)
				}

				if i, ok := test.court["Name"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Name\"]' should be a string")
					}
					court.Container.Name = value
					assert.Equal(t, court.Container.Name, value, "The Court name was not updated correctly")
				}

				if i, ok := test.court["Players"]; ok {
					value, ok := i.([]string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Players\"]' should be an array of strings")
					}
					court.Container.Players = value
					if !utilities.Equal(court.Container.Players, value) {
						t.Errorf("The Court name was not updated correctly:\n got %v\n want %v", court.Container.Players, value)
					}
				}
			}
		})
	}
}
