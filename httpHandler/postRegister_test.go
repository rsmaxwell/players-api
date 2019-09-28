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

func TestRegister(t *testing.T) {

	teardown := SetupEmpty(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		userID         string
		password       string
		firstName      string
		lastName       string
		email          string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			userID:         "007",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Space in userID",
			userID:         "0 7",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "Path in userID",
			userID:         "../007",
			password:       "topsecret",
			firstName:      "James",
			lastName:       "Bond",
			email:          "james@mi6.co.uk",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(RegisterRequest{
				UserID:    test.userID,
				Password:  test.password,
				FirstName: test.firstName,
				LastName:  test.lastName,
				Email:     test.email,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rw := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can call its ServeHTTP method
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
				assert.Equal(t, initialNumberOfPeople+1, finalNumberOfPeople, "Person was not registered")
			} else {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople, "Unexpected number of people")
			}
		})
	}
}
