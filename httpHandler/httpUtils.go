package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/logger"
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
	logger.Logger.Printf("Response: %d %s - %s", httpStatus, http.StatusText(httpStatus), message)
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

// Simple check on the user calling the service
func checkUser(userID, password string) bool {

	person, err := person.GetPersonDetails(userID)
	if err != nil {
		logger.Logger.Printf("checkUser(0): userID or password do not match")
		return false
	}

	err = bcrypt.CompareHashAndPassword(person.HashedPassword, []byte(password))
	if err != nil {
		logger.Logger.Printf("checkUser(1): userID or password do not match")
		return false
	}

	logger.Logger.Printf("checkUser: OK\n")
	return true
}