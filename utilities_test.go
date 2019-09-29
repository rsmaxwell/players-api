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
	"testing"

	"github.com/rsmaxwell/players-api/destination"
	"github.com/rsmaxwell/players-api/httphandler"
	"github.com/rsmaxwell/players-api/person"
)

var (
	CurrentUserID   = "007"
	CurrentPassword = "topsecret"
	CurrentToken    string
	AnotherUserID   = "bob"
	CurrentCourtID  = "1000"
	AllPeopleIDs    []string
	AllCourtIDs     []string
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

	httphandler.Register(rw, req)
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

	httphandler.Login(rw, req)

	bytes, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response httphandler.LogonResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		log.Fatalln(err)
	}

	return response.Token
}

func UpdatePerson(token, userID string, person2 map[string]interface{}) {

	request := httphandler.UpdatePersonRequest{Token: token, Person: person2}
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	rw := httptest.NewRecorder()

	httphandler.UpdatePerson(rw, req, userID)

	if rw.Code != http.StatusOK {
		log.Fatalln(err)
	}
}

func CreateCourt(token, name string, players []string) {
	requestBody, err := json.Marshal(httphandler.CreateCourtRequest{
		Token: token,
		Court: destination.Court{
			Container: destination.PeopleContainer{
				Name:    name,
				Players: players,
			},
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

	httphandler.CreateCourt(rw, req)

	if rw.Code != 200 {
		bytes, err := ioutil.ReadAll(rw.Body)
		if err != nil {
			log.Fatalln(err)
		}
		response := string(bytes)
		log.Fatalln(fmt.Printf("Code: %d, Response: %s", rw.Code, response))
	}
}

func setupEmpty(t *testing.T) func(t *testing.T) {
	person.Clear()
	destination.ClearCourts()
	destination.ClearQueue()

	return func(t *testing.T) {
		person.Clear()
		destination.ClearCourts()
		destination.ClearQueue()
	}
}

func setupOne(t *testing.T) func(t *testing.T) {
	person.Clear()
	destination.ClearCourts()
	destination.ClearQueue()

	Register("007", "topsecret", "James", "Bond", "james@mi6.co.uk")

	AllPeopleIDs = []string{}
	AllCourtIDs = []string{}

	return func(t *testing.T) {
		person.Clear()
		destination.ClearCourts()
		destination.ClearQueue()
	}
}

func setupLoggedin(t *testing.T) func(t *testing.T) {
	person.Clear()
	destination.ClearCourts()
	destination.ClearQueue()

	Register("007", "topsecret", "James", "Bond", "james@mi6.co.uk")

	AllPeopleIDs = []string{"007"}
	AllCourtIDs = []string{}

	CurrentToken = Login(CurrentUserID, CurrentPassword)

	return func(t *testing.T) {
		person.Clear()
		destination.ClearCourts()
		destination.ClearQueue()
	}
}

func setupFull(t *testing.T) func(t *testing.T) {

	person.Clear()
	destination.ClearCourts()
	destination.ClearQueue()

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

	CurrentToken = Login(CurrentUserID, CurrentPassword)

	// Make all the people a 'player'
	person2 := make(map[string]interface{})
	person2["Player"] = true

	AllPeopleIDs = []string{"007", "bob", "alice", "jill", "david", "mary", "john", "judith", "paul", "nigel", "jeremy", "joanna"}
	for _, id := range AllPeopleIDs {
		UpdatePerson(CurrentToken, id, person2)
	}

	CreateCourt(CurrentToken, "Court 1", []string{})
	CreateCourt(CurrentToken, "Court 2", []string{})

	AllCourtIDs = []string{"1000", "1001"}

	return func(t *testing.T) {
		person.Clear()
		destination.ClearCourts()
		destination.ClearQueue()
	}
}
