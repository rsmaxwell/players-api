package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestNotFound(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			token:          MyToken,
			expectedStatus: http.StatusNotFound,
			expectedResult: AllPeopleIDs,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(ListPeopleRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a  request to pass to our handler.
			req, err := http.NewRequest("GET", "/junk", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			SetupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, test.expectedStatus)
			}
		})
	}
}
