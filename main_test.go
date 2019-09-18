package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/httpHandler"
	"github.com/rsmaxwell/players-api/person"
	"github.com/stretchr/testify/assert"
)

func TestGetPerson(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		userID         string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			userID:         CurrentUserID,
			expectedStatus: http.StatusOK,
			expectedResult: "James",
		},
		{
			name:           "Bad token",
			token:          "junk",
			userID:         CurrentUserID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			name:           "Bad userID",
			token:          CurrentToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.GetPersonRequest{
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

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var response httpHandler.GetPersonResponse

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

func TestListPeople(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			expectedStatus: http.StatusOK,
			expectedResult: AllPeopleIDs,
		},
		{
			name:           "Bad token",
			token:          "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.ListPeopleRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a  request to pass to our handler.
			req, err := http.NewRequest("GET", "/person", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

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

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			if rr.Code == http.StatusOK {
				var response httpHandler.ListPeopleResponse

				err = json.Unmarshal(bytes, &response)
				if err != nil {
					log.Fatalln(err)
				}

				// Check the response body is what we expect.
				if !Equal(response.People, test.expectedResult) {
					t.Logf("actual:   %v", response.People)
					t.Logf("expected: %v", test.expectedResult)
					t.Errorf("Unexpected list of people")
				}
			}
		})
	}
}

func TestDeletePerson(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			userID:         AnotherUserID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Bad token",
			token:          "junk",
			userID:         AnotherUserID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Bad userID",
			token:          CurrentToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			initialNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(httpHandler.ListPeopleRequest{
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

			finalNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			if rr.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople+1, "Person was not deleted")
			} else {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople, "Unexpected number of people")
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name            string
		token           string
		expectedStatus  int
		expectedResults int
	}{
		{
			name:            "Good request",
			token:           CurrentToken,
			expectedStatus:  http.StatusOK,
			expectedResults: 0,
		},
		{
			name:            "Bad token",
			token:           "junk",
			expectedStatus:  http.StatusUnauthorized,
			expectedResults: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.GetMetricsRequest{
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

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			if rr.Code == http.StatusOK {
				var response httpHandler.GetMetricsResponse

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

// -----------------------------------------------------------------------------------------------------

func TestGetCourt(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		courtID        string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusOK,
			expectedResult: "Court 1",
		},
		{
			name:           "Bad token",
			token:          "junk",
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			name:           "Bad userID",
			token:          CurrentToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.GetCourtRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/court/"+test.courtID, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

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

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var response httpHandler.GetCourtResponse

			err = json.Unmarshal(bytes, &response)
			if err != nil {
				log.Fatalln(err)
			}

			// Check the response body is what we expect.
			actual := response.Court.Name
			if actual != test.expectedResult {
				t.Logf("actual:   <%v>", actual)
				t.Logf("expected: <%v>", test.expectedResult)
				t.Errorf("handler returned unexpected body: got %v want %v", actual, test.expectedResult)
			}
		})
	}
}

func TestListCourts(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			expectedStatus: http.StatusOK,
			expectedResult: AllCourtIDs,
		},
		{
			name:           "Bad token",
			token:          "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.ListCourtsRequest{
				Token: test.token,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a  request to pass to our handler.
			req, err := http.NewRequest("GET", "/court", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

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

			bytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				log.Fatalln(err)
			}

			if rr.Code == http.StatusOK {
				var response httpHandler.ListCourtsResponse

				err = json.Unmarshal(bytes, &response)
				if err != nil {
					log.Fatalln(err)
				}

				// Check the response body is what we expect.
				if !Equal(response.Courts, test.expectedResult) {
					t.Logf("actual:   %v", response.Courts)
					t.Logf("expected: %v", test.expectedResult)
					t.Errorf("Unexpected list of courts")
				}
			}
		})
	}
}

func TestDeleteCourt(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		name           string
		token          string
		courtID        string
		expectedStatus int
	}{
		{
			name:           "Good request",
			token:          CurrentToken,
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Bad token",
			token:          "junk",
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Bad userID",
			token:          CurrentToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			initialNumberOfCourts, err := court.Size()
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(httpHandler.ListCourtsRequest{
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

			finalNumberOfCourts, err := court.Size()
			if err != nil {
				t.Fatal(err)
			}

			if rr.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfCourts, finalNumberOfCourts+1, "Court was not deleted")
			} else {
				assert.Equal(t, initialNumberOfCourts, finalNumberOfCourts, "Unexpected number of courts")
			}
		})
	}
}
