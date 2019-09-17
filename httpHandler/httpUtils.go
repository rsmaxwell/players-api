package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
	"golang.org/x/crypto/bcrypt"
)

var (
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// Error Response
type messageResponse struct {
	Message string `json:"message"`
}

// WriteResponse method
func WriteResponse(w http.ResponseWriter, httpStatus int, message string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	json.NewEncoder(w).Encode(messageResponse{
		Message: message,
	})
}

func setHeaders(rw http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	if origin == "" {
		origin = "http://localhost:4200"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", origin)
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Authorization")
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
