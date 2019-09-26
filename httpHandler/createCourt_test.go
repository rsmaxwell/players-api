package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/container"
	"github.com/rsmaxwell/players-api/court"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateCourt(t *testing.T) {

	teardown := SetupLoggedin(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		name           string
		players        []string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          MyToken,
			name:           "Court 1",
			players:        []string{},
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			name:           "Court 1",
			players:        []string{},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfCourts, err := court.Size()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(CreateCourtRequest{
				Token: test.token,
				Court: court.Court{
					Container: container.Container{
						Name:    test.name,
						Players: test.players,
					},
				},
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("POST", "/court", bytes.NewBuffer(requestBody))
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

			finalNumberOfCourts, err := court.Size()
			if err != nil {
				t.Fatal(err)
			}

			if rr.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfCourts+1, finalNumberOfCourts, "Court was not registered")
			} else {
				assert.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
