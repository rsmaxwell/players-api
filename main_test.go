package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsmaxwell/players-api/session"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/httpHandler"
	"github.com/rsmaxwell/players-api/person"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	teardown := setupEmpty(t)
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

			requestBody, err := json.Marshal(httpHandler.RegisterRequest{
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

			finalNumberOfPeople, err := person.Size()
			if err != nil {
				t.Fatal(err)
			}

			if rr.Code == http.StatusOK {
				assert.Equal(t, initialNumberOfPeople+1, finalNumberOfPeople, "Person was not registered")
			} else {
				assert.Equal(t, initialNumberOfPeople, finalNumberOfPeople, "Unexpected number of people")
			}
		})
	}
}

func TestLogin(t *testing.T) {

	teardown := setupOne(t)
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

			req.Header.Set("Authorization", "Basic "+basicAuth(test.userID, test.password))

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

			if rr.Code == http.StatusOK {
				// Check the response is what we expect.
				bytes, err := ioutil.ReadAll(rr.Body)
				if err != nil {
					log.Fatalln(err)
				}

				var response httpHandler.LogonResponse

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

func TestGetPerson(t *testing.T) {

	teardown := setupFull(t)
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
			token:          CurrentToken,
			userID:         CurrentUserID,
			expectedStatus: http.StatusOK,
			expectedResult: "James",
		},
		{
			testName:       "Bad token",
			token:          "junk",
			userID:         CurrentUserID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			testName:       "Bad userID",
			token:          CurrentToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
		testName       string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			token:          CurrentToken,
			expectedStatus: http.StatusOK,
			expectedResult: AllPeopleIDs,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
		testName       string
		token          string
		userID         string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          CurrentToken,
			userID:         AnotherUserID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			userID:         AnotherUserID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			token:          CurrentToken,
			userID:         "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
		testName        string
		token           string
		expectedStatus  int
		expectedResults int
	}{
		{
			testName:        "Good request",
			token:           CurrentToken,
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

func TestCreateCourt(t *testing.T) {

	teardown := setupLoggedin(t)
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
			token:          CurrentToken,
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

			requestBody, err := json.Marshal(httpHandler.CreateCourtRequest{
				Token: test.token,
				Court: court.Court{
					Name:    test.name,
					Players: test.players,
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

func TestUpdateCourt(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		id             string
		court          map[string]interface{}
		expectedStatus int
	}{
		{
			testName: "Good request",
			token:    CurrentToken,
			id:       CurrentCourtID,
			court: map[string]interface{}{
				"Name":    "COURT 101",
				"Players": []string{"bob", "jill", "alice"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Bad token",
			token:    "junk",
			id:       CurrentCourtID,
			court: map[string]interface{}{
				"Name":    "COURT 101",
				"Players": []string{},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName: "Bad userID",
			token:    CurrentToken,
			id:       "junk",
			court: map[string]interface{}{
				"Name":    "COURT 101",
				"Players": []string{},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			testName: "Bad player",
			token:    CurrentToken,
			id:       CurrentCourtID,
			court: map[string]interface{}{
				"Name":    "COURT 101",
				"Players": []string{"junk"},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, err := json.Marshal(httpHandler.UpdateCourtRequest{
				Token: test.token,
				Court: test.court,
			})
			if err != nil {
				log.Fatalln(err)
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest("PUT", "/court/"+test.id, bytes.NewBuffer(requestBody))
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

			// Check the court was actually updated
			if rr.Code == http.StatusOK {
				court, err := court.Load(test.id)
				if err != nil {
					t.Fatal(err)
				}

				if i, ok := test.court["Name"]; ok {
					value, ok := i.(string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Name\"]' should be a string")
					}
					court.Name = value
					assert.Equal(t, court.Name, value, "The Court name was not updated correctly")
				}

				if i, ok := test.court["Players"]; ok {
					value, ok := i.([]string)
					if !ok {
						t.Errorf("The type of 'test.court[\"Players\"]' should be an array of strings")
					}
					court.Players = value
					if !Equal(court.Players, value) {
						t.Errorf("The Court name was not updated correctly:\n got %v\n want %v", court.Players, value)
					}
				}
			}
		})
	}
}

func TestGetCourt(t *testing.T) {

	teardown := setupFull(t)
	defer teardown(t)

	tests := []struct {
		testName       string
		token          string
		courtID        string
		expectedStatus int
		expectedResult string
	}{
		{
			testName:       "Good request",
			token:          CurrentToken,
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusOK,
			expectedResult: "Court 1",
		},
		{
			testName:       "Bad token",
			token:          "junk",
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusUnauthorized,
			expectedResult: "",
		},
		{
			testName:       "Bad userID",
			token:          CurrentToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
		testName       string
		token          string
		expectedStatus int
		expectedResult []string
	}{
		{
			testName:       "Good request",
			token:          CurrentToken,
			expectedStatus: http.StatusOK,
			expectedResult: AllCourtIDs,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
		testName       string
		token          string
		courtID        string
		expectedStatus int
	}{
		{
			testName:       "Good request",
			token:          CurrentToken,
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "Bad token",
			token:          "junk",
			courtID:        CurrentCourtID,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName:       "Bad userID",
			token:          CurrentToken,
			courtID:        "junk",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

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
