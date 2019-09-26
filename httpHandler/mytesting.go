package httphandler

import (
	"encoding/base64"
	"testing"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/queue"
	"github.com/rsmaxwell/players-api/session"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	// MySession logged on session
	MySession *session.Session

	// MyToken logged on user
	MyToken string

	// MyUserID logged on user
	MyUserID = "007"

	// AnotherUserID for testing
	AnotherUserID = "bob"

	// MyCourtID is a court created at setup
	MyCourtID = "1000"

	// AllPeopleIDs is a list of the people created at setup
	AllPeopleIDs []string

	// AllCourtIDs is a list of the courts created at setup
	AllCourtIDs []string
)

// BasicAuth converts a username and password into BasicAuth format
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func registerPerson(r *require.Assertions, id, password, firstname, lastname, email string) {

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	r.Nil(err, "err should be nothing")

	// Make the person
	err = person.Save(id, person.New(firstname, lastname, email, hashedPassword, false))
	r.Nil(err, "err should be nothing")

	AllPeopleIDs = append(AllPeopleIDs, id)
}

func createPlayer(r *require.Assertions, id, password, firstname, lastname, email string) {

	// Register a person
	registerPerson(r, id, password, firstname, lastname, email)

	// Make the Person a 'player'
	person2 := make(map[string]interface{})
	person2["Player"] = true

	_, err := person.Update(id, MySession, person2)
	r.Nil(err, "err should be nothing")
}

func getSession(r *require.Assertions, id string) {

	var err error

	MyToken, err = session.New(id)
	r.Nil(err, "err should be nothing")

	MySession = session.LookupToken(MyToken)
}

func createplayers(r *require.Assertions) {

	AllPeopleIDs = []string{}

	createPlayer(r, "007", "topsecret", "James", "Bond", "james@mi6.uk.gov")
	createPlayer(r, "bob", "qwerty", "Robert", "Bruce", "bob@aol.com")
	createPlayer(r, "alice", "wonder", "Alice", "Wonderland", "alice@abc.com")
	createPlayer(r, "jill", "password", "Jill", "Cooper", "jill@def.com")
	createPlayer(r, "david", "magic", "David", "Copperfield", "david@ghi.com")
	createPlayer(r, "mary", "queen", "Mary", "Gray", "mary@jkl.com")
	createPlayer(r, "john", "king", "John", "King", "john@mno.com")
	createPlayer(r, "judith", "bean", "Judith", "Green", "judith@pqr.com")
	createPlayer(r, "paul", "ruler", "Paul", "Straight", "paul@stu.com")
	createPlayer(r, "nigel", "careful", "Nigel", "Curver", "nigel@vwx.com")
	createPlayer(r, "jeremy", "changeme", "Jeremy", "Black", "jeremy@vwx.com")
	createPlayer(r, "joanna", "bright", "Joanna", "Brown", "joanna@yza.com")
}

func createcourt(r *require.Assertions, name string) {

	id, err := court.Insert(court.New(name))
	r.Nil(err, "err should be nothing")

	AllCourtIDs = append(AllCourtIDs, id)
}

func createcourts(r *require.Assertions) {

	AllCourtIDs = []string{}

	createcourt(r, "Court 1")
	createcourt(r, "Court 2")
}

// SetupEmpty creates a single person
func SetupEmpty(t *testing.T) func(t *testing.T) {
	person.Clear()
	court.Clear()
	queue.Clear()

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
		queue.Clear()
	}
}

// SetupOne creates a single person
func SetupOne(t *testing.T) func(t *testing.T) {
	r := require.New(t)

	person.Clear()
	court.Clear()
	queue.Clear()

	AllPeopleIDs = []string{}
	AllCourtIDs = []string{}

	registerPerson(r, "007", "topsecret", "James", "Bond", "james@mi6.uk.gov")

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
		queue.Clear()
	}
}

// SetupLoggedin creates a logged in player
func SetupLoggedin(t *testing.T) func(t *testing.T) {
	r := require.New(t)

	person.Clear()
	court.Clear()
	queue.Clear()

	getSession(r, "007")

	AllPeopleIDs = []string{}
	AllCourtIDs = []string{}

	createPlayer(r, "007", "topsecret", "James", "Bond", "james@mi6.uk.gov")

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
		queue.Clear()
	}
}

// SetupFull creates a logged in players and a number of other people and courts
func SetupFull(t *testing.T) func(t *testing.T) {
	r := require.New(t)

	person.Clear()
	court.Clear()
	queue.Clear()

	getSession(r, "007")

	createplayers(r)
	createcourts(r)

	startup()

	return func(t *testing.T) {
		person.Clear()
		court.Clear()
		queue.Clear()
	}
}
