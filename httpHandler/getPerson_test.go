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

func TestGetPerson(t *testing.T) {

	teardown := SetupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		userID         string
		expectedStatus int
		expectedResult string
	}{
		{
			testName:       "Good request",
			token:          MyToken,
			userID:         MyUserID,
			expectedStatus: http.StatusOK,
			expectedResult: "James",
		},
		{
			testName:       "Bad token",
			token:          "junk",
			userID:         MyUserID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			testName:       "Bad userID",
			token:          MyToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(GetPersonRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/person/"+test.userID, bytes.NewBuffer(requestBody))
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
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var response GetPersonResponse

			err = json.Unmarshal(bytes, &response)
			if err != nil {
				log.Fatalln(err)
			}

			// Check the response body is what we expect.
			actual := response.Person.FirstName
			if actual != test.expectedResult {
				t.Logf("actual:   <%v>", actual)
				t.Logf("expected: <%v>", test.expectedResult)
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}
