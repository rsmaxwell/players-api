package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/session"

	"github.com/gorilla/mux"
)

func TestLogin(t *testing.T) {

	teardown := SetupOne(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		userID         string
		password       string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			userID:         "007",
			password:       "topsecret",
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Space in userID",
			userID:         "0 7",
			password:       "topsecret",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Path in userID",
			userID:         "../007",
			password:       "topsecret",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/login", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", BasicAuth(test.userID, test.password))

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

			if rr.Code == http.StatusOK {
				// Check the response is what we expect.
				bytes, err := ioutil.ReadAll(rr.Body)
				if err != nil {
					log.Fatalln(err)
				}

				var response LogonResponse

				err = json.Unmarshal(bytes, &response)
				if err != nil {
					log.Fatalln(err)
				}

				// Check the response contains a valid token.
				session := session.LookupToken(response.Token)
				if session == nil {
					t.Errorf("Invalid token returned: token:%s", response.Token)
				}
			}
		})
	}
}
