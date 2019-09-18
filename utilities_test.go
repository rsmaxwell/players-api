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

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/httpHandler"
	"github.com/rsmaxwell/players-api/logger"
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

func setupEmpty(t *testing.T) func(t *testing.T) {
	person.Clear()
	court.Clear()

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
	}
}

func setupFull(t *testing.T) func(t *testing.T) {

	person.Clear()
	court.Clear()

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

	json := "{\"player\":true}"
	AllPeopleIDs = []string{"007", "bob", "alice", "jill", "david", "mary", "john", "judith", "paul", "nigel", "jeremy", "joanna"}
	for _, id := range AllPeopleIDs {
		UpdatePerson(CurrentToken, id, json)
	}

	CreateCourt(CurrentToken, "Court 1", []string{})
	CreateCourt(CurrentToken, "Court 2", []string{})

	AllCourtIDs = []string{"1000", "1001"}

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
	}
}

// Equal tells whether a and b contain the same elements without order
func Equal(x, y []string) bool {

	if x == nil {
		if y == nil {
			return true
		}
		return false
	} else if y == nil {
		return false
	}

	if len(x) != len(y) {
		return false
	}

	xMap := make(map[string]int)
	yMap := make(map[string]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}
