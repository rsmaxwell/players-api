package httphandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetMetrics(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName        string
		token           string
		expectedStatus  int
		expectedResults int
	}{
		{
			testName:        "Good request",
			token:           MyToken,
			expectedStatus:  http.StatusOK,
			expectedResults: 0,
		},
		{
			testName:        "Bad token",
			token:           "junk",
			expectedStatus:  http.StatusUnauthorized,
			expectedResults: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetMetricsRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/metrics", bytes.NewBuffer(requestBody))
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

			bytes, err := ioutil.ReadAll(rw.Body)
			if err != nil {
				log.Fatalln(err)
			}

			if rw.Code == http.StatusOK {
				var response GetMetricsResponse

				err = json.Unmarshal(bytes, &response)
				if err != nil {
					log.Fatalln(err)
				}

				// Check the response body is what we expect.
				if response.ClientSuccess != test.expectedResults {
					t.Logf("actual:   %v", response.ClientSuccess)
					t.Logf("expected: %v", test.expectedResults)
					t.Errorf("Unexpected metrics")
				}
			}
		})
	}
}
