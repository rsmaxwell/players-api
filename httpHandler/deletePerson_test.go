package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/person"
	"github.com/stretchr/testify/assert"
)

func TestDeletePerson(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		userID         string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          MyToken,
			userID:         AnotherUserID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			userID:         AnotherUserID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			token:          MyToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(ListPeopleRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("DELETE", "/person/"+test.userID, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rw := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rw, req)

			// Check the status code is what we expect.
			if rw.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rw.Code, test.expectedStatus)
			}

			finalNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			if rw.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople+1, "Person was not deleted")
			} else {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople, "Unexpected number of people")
			}
		})
	}
}
