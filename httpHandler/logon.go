package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/session"
	"golang.org/x/crypto/bcrypt"
)

// LogonResponse structure
type LogonResponse struct {
	Token string `json:"token"`
}

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !CheckUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	token, err := session.New(user)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, "Error creating session")
		serverError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(LogonResponse{
		Token: token,
	})
}

// CheckUser - Basic check on the user calling the service
func CheckUser(id, password string) bool {

	person, err := person.Get(id)
	if err != nil {
		return false
	}
	if person == nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(person.HashedPassword, []byte(password))
	if err != nil {
		return false
	}

	return true
}
