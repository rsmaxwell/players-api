package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/players"
	"github.com/stretchr/testify/assert"
)

// assert that the response is JSON of the expected format
func assertResponseJSON(t *testing.T, rr *httptest.ResponseRecorder, expectedObj interface{}) {
	rr.Flush()
	assert.Equal(t, rr.Header()["Content-Type"], []string{"application/json"}, "Unexpected Content-Type")
	ej, err := json.Marshal(expectedObj)
	assert.Nil(t, err)
	assert.JSONEq(t, string(ej), string(rr.Body.Bytes()))
}

func TestCheckuser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		expectedResult bool
	}{
		{
			name:           "Good credentials",
			username:       "fred",
			password:       "bloggs",
			expectedResult: true,
		},
		{
			name:           "Bad credentials",
			username:       "one",
			password:       "two",
			expectedResult: false,
		},
	}

	username = "fred"
	password = "bloggs"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := checkUser(test.username, test.password)
			assert.Equal(t, test.expectedResult, ok)
		})
	}
}

// Testing the writeStatusResponse Function
func TestWriteMessageResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	msg := "A message of sorts"
	status := http.StatusInternalServerError

	expectedObj := messageResponseJSON{
		Message: msg,
	}

	writeMessageResponse(rr, status, msg)
	assert.Equal(t, status, rr.Code, "Wrong HTTP status code returned")
	assertResponseJSON(t, rr, expectedObj)
}

func Reset(names ...string) error {

	err := players.Reset(names...)
	if err != nil {
		return err
	}

	username = "fred"
	password = "bloggs"

	clientSuccess = 0
	clientError = 0
	clientAuthenticationError = 0
	serverError = 0

	return nil
}

func TestWritePlayerInfoResponse(t *testing.T) {

	err := Reset("fred")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			username:       username,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"numberOfPeople":1}`,
		},
		{
			name:           "Bad credentials",
			username:       "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: `{"message":"Invalid username and/or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/personinfo", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can call its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Logf("actual:   <%v>", actual)
				t.Logf("expected: <%v>", test.expectedResult)
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}

func TestWriteGetListOfPlayersResponse(t *testing.T) {

	err := Reset()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			username:       username,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"people":null}`,
		},
		{
			name:           "Bad credentials",
			username:       "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: `{"message":"Invalid username and/or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a  request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/people", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}

func TestWriteGetPlayerDetailsResponse(t *testing.T) {

	err := Reset("fred", "bloggs")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			username:       username,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"person":{"name":"fred"}}`,
		},
		{
			name:           "Bad credentials",
			username:       "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: `{"message":"Invalid username and/or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Create a  request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/person/1000", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}

func TestAddPlayerResponse(t *testing.T) {

	err := Reset("fred")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		username               string
		password               string
		expectedStatus         int
		expectedResult         string
		expectedNumberOfPeople int
	}{
		{
			name:                   "Good credentials",
			username:               username,
			password:               password,
			expectedStatus:         http.StatusOK,
			expectedResult:         `{"message":"ok"}`,
			expectedNumberOfPeople: 2,
		},
		{
			name:                   "Bad credentials",
			username:               "junk",
			password:               "junk",
			expectedStatus:         http.StatusUnauthorized,
			expectedResult:         `{"message":"Invalid username and/or password"}`,
			expectedNumberOfPeople: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Create a  request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/person", strings.NewReader("{\"name\":\"Richard\"}"))
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}

			list, err := players.List()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(list), test.expectedNumberOfPeople, "Unexpected number of Players: got %v want %v", actual, test.expectedNumberOfPeople)
		})
	}
}

func TestDeletePlayerResponse(t *testing.T) {

	err := Reset("fred", "bloggs")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		username               string
		password               string
		expectedStatus         int
		expectedResult         string
		expectedNumberOfPeople int
	}{
		{
			name:                   "Good credentials",
			username:               username,
			password:               password,
			expectedStatus:         http.StatusOK,
			expectedResult:         `{"message":"ok"}`,
			expectedNumberOfPeople: 1,
		},
		{
			name:                   "Bad credentials",
			username:               "junk",
			password:               "junk",
			expectedStatus:         http.StatusUnauthorized,
			expectedResult:         `{"message":"Invalid username and/or password"}`,
			expectedNumberOfPeople: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Create a  request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", "/person/1000", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}

			list, err := players.List()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Number of Players = %d", len(list))
			assert.Equal(t, len(list), test.expectedNumberOfPeople, "Unexpected number of People: got %v want %v", len(list), test.expectedNumberOfPeople)
		})
	}
}

func TestWriteGetMetricsResponse(t *testing.T) {

	err := Reset("fred", "bloggs")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			username:       username,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"clientSuccess":0,"clientError":0,"clientAuthenticationError":0,"serverError":0}`,
		},
		{
			name:           "Bad credentials",
			username:       "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: `{"message":"Invalid username and/or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a  request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/metrics", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.username, test.password)

			router := mux.NewRouter()
			setupHandlers(router)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our router satisfies http.Handler, so we can its ServeHTTP method
			// directly and pass in our ResponseRecorder and Request.
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
			}

			// Check the response body is what we expect.
			actual := strings.TrimRight(rr.Body.String(), "\n")
			if actual != test.expectedResult {
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}
