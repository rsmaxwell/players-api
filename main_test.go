package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/httpHandler"
	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/person"
	"github.com/stretchr/testify/assert"
)

var (
	userID   = "007"
	password = "topsecret"
	token    string
)

func Register(userID, password, firstName, lastName, email string) {
	requestBody, err := json.Marshal(map[string]string{
		"UserID":    userID,
		"Password":  password,
		"FirstName": firstName,
		"LastName":  lastName,
		"Email":     email,
	})
	if err != nil {
		log.Fatalln(err)
	}

	rw := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	httpHandler.Register(rw, req)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Login(userID, password string) string {

	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(userID, password))

	rw := httptest.NewRecorder()

	httpHandler.Login(rw, req)

	bytes, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response httpHandler.AuthenticateResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		log.Fatalln(err)
	}

	return response.Token
}

func UpdatePerson(token, userID, json string) {

	requestBody := []byte("{ \"token\": \"" + token + "\", \"person\":" + json + "}")

	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	rw := httptest.NewRecorder()

	httpHandler.UpdatePerson(rw, req, userID)
}

func CreateCourt(token, name string, players []string) {
	requestBody, err := json.Marshal(httpHandler.CreateCourtRequest{
		Token: token,
		Court: court.Court{
			Name:    name,
			Players: players,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	rw := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	httpHandler.CreateCourt(rw, req)

	if rw.Code != 200 {
		bytes, err := ioutil.ReadAll(rw.Body)
		if err != nil {
			log.Fatalln(err)
		}
		response := string(bytes)
		logger.Logger.Fatalln(fmt.Printf("Code: %d, Response: %s", rw.Code, response))
	}
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	logger.Logger.Printf("setup test case")

	logger.Logger.Printf("Clear the people and courts")
	person.Clear()
	court.Clear()

	logger.Logger.Printf("Create a list of people")
	Register("007", "topsecret", "James", "Bond", "james@mi6.co.uk")
	Register("bob", "qwerty", "Robert", "Bruce", "bob@aol.com")
	Register("alice", "wonder", "Alice", "Wonderland", "alice@abc.com")
	Register("jill", "password", "Jill", "Cooper", "jill@def.com")
	Register("david", "magic", "David", "Copperfield", "david@ghi.com")
	Register("mary", "queen", "Mary", "Gray", "mary@jkl.com")
	Register("john", "king", "John", "King", "john@mno.com")
	Register("judith", "bean", "Judith", "Green", "judith@pqr.com")
	Register("paul", "ruler", "Paul", "Straight", "paul@stu.com")
	Register("nigel", "careful", "Nigel", "Curver", "nigel@vwx.com")
	Register("jeremy", "changeme", "Jeremy", "Black", "jeremy@vwx.com")
	Register("joanna", "bright", "Joanna", "Brown", "joanna@yza.com")

	logger.Logger.Printf("Login as [\"007\"]")
	token = Login("007", "topsecret")
	logger.Logger.Printf("token: %s", token)

	logger.Logger.Printf("Update people (to make them players)")
	json := "{\"player\":true}"
	ids := []string{"007", "bob", "alice", "jill", "david", "mary", "john", "judith", "paul", "nigel", "jeremy", "joanna"}
	for _, id := range ids {
		UpdatePerson(token, id, json)
	}

	logger.Logger.Printf("Create a list of courts")
	CreateCourt(token, "Court 1", []string{})
	CreateCourt(token, "Court 2", []string{})

	return func(t *testing.T) {
		t.Log("teardown test case")
		person.Clear()
		court.Clear()
	}
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
			username:       "007",
			password:       "topsecret",
			expectedResult: true,
		},
		{
			name:           "Bad credentials",
			username:       "007",
			password:       "junk",
			expectedResult: false,
		},
		{
			name:           "non-existant user",
			username:       "junk",
			password:       "junk",
			expectedResult: false,
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := httpHandler.CheckUser(test.username, test.password)
			assert.Equal(t, test.expectedResult, ok)
		})
	}
}

func TestWritePlayerInfoResponse(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name           string
		userID         string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			userID:         userID,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"numberOfPeople":1}`,
		},
		{
			name:           "Bad credentials",
			userID:         "junk",
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
			req.SetBasicAuth(test.userID, test.password)

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

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name           string
		userID         string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			userID:         userID,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"people":null}`,
		},
		{
			name:           "Bad credentials",
			userID:         "junk",
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
			req.SetBasicAuth(test.userID, test.password)

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

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name           string
		userID         string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			userID:         userID,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"person":{"name":"fred"}}`,
		},
		{
			name:           "Bad credentials",
			userID:         "junk",
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
			req.SetBasicAuth(test.userID, test.password)

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

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name                   string
		userID                 string
		password               string
		expectedStatus         int
		expectedResult         string
		expectedNumberOfPeople int
	}{
		{
			name:                   "Good credentials",
			userID:                 userID,
			password:               password,
			expectedStatus:         http.StatusOK,
			expectedResult:         `{"message":"ok"}`,
			expectedNumberOfPeople: 2,
		},
		{
			name:                   "Bad credentials",
			userID:                 "junk",
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
			req.SetBasicAuth(test.userID, test.password)

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

			list, err := person.List()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(list), test.expectedNumberOfPeople, "Unexpected number of Players: got %v want %v", actual, test.expectedNumberOfPeople)
		})
	}
}

func TestDeletePlayerResponse(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name                   string
		userID                 string
		password               string
		expectedStatus         int
		expectedResult         string
		expectedNumberOfPeople int
	}{
		{
			name:                   "Good credentials",
			userID:                 userID,
			password:               password,
			expectedStatus:         http.StatusOK,
			expectedResult:         `{"message":"ok"}`,
			expectedNumberOfPeople: 1,
		},
		{
			name:                   "Bad credentials",
			userID:                 "junk",
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
			req.SetBasicAuth(test.userID, test.password)

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

			list, err := person.List()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Number of Players = %d", len(list))
			assert.Equal(t, len(list), test.expectedNumberOfPeople, "Unexpected number of People: got %v want %v", len(list), test.expectedNumberOfPeople)
		})
	}
}

func TestWriteGetMetricsResponse(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name           string
		userID         string
		password       string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Good credentials",
			userID:         userID,
			password:       password,
			expectedStatus: http.StatusOK,
			expectedResult: `{"clientSuccess":0,"clientError":0,"clientAuthenticationError":0,"serverError":0}`,
		},
		{
			name:           "Bad credentials",
			userID:         "junk",
			password:       "junk",
			expectedStatus: http.StatusUnauthorized,
			expectedResult: `{"message":"Invalid username and/or password"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/metrics", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(test.userID, test.password)

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
