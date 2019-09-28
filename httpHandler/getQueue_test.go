package httphandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/utilities"

	"github.com/gorilla/mux"
)

func TestGetQueue(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName              string
		token                 string
		userID                string
		expectedStatus        int
		expectedResultName    string
		expectedResultPlayers []string
	}{
		{
			testName:              "Good request",
			token:                 MyToken,
			expectedStatus:        http.StatusOK,
			expectedResultName:    "James",
			expectedResultPlayers: []string{"one", "two"},
		},
		{
			testName:              "Bad token",
			token:                 "junk",
			expectedStatus:        http.StatusUnauthorized,
			expectedResultName:    "",
			expectedResultPlayers: []string{"one", "two"},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetQueueRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/queue", bytes.NewBuffer(requestBody))
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

			bytes, err := ioutil.ReadAll(rw.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var response GetQueueResponse

			err = json.Unmarshal(bytes, &response)
			if err != nil {
				log.Fatalln(err)
			}

			if rw.Code != http.StatusOK {

				// Check the response body is what we expect.
				actualName := response.Queue.Container.Name
				if actualName != test.expectedResultName {
					t.Logf("actual:   <%v>", actualName)
					t.Logf("expected: <%v>", test.expectedResultName)
					t.Errorf("handler returned unexpected body: got %v want %v", actualName, test.expectedResultName)
				}

				actualPlayers := response.Queue.Container.Players
				if utilities.Equal(actualPlayers, test.expectedResultPlayers) {
					t.Logf("actual:   <%v>", actualPlayers)
					t.Logf("expected: <%v>", test.expectedResultPlayers)
					t.Errorf("handler returned unexpected body: got %v want %v", actualPlayers, test.expectedResultPlayers)
				}
			}
		})
	}
}
