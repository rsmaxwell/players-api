package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/destination"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCourt(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		courtID        string
		expectedStatus int
	}{
		{
			testName:       "Good  request",
			token:          MyToken,
			courtID:        MyCourtID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			courtID:        MyCourtID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			token:          MyToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			initialNumberOfCourts, err := destination.CourtSize()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(ListCourtsRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("DELETE", "/court/"+test.courtID, bytes.NewBuffer(requestBody))
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

			finalNumberOfCourts, err := destination.CourtSize()
			if err != nil {
				t.Fatal(err)
			}

			if rw.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfCourts, finalNumberOfCourts+1, "Court was not deleted")
			} else {
				assert.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
