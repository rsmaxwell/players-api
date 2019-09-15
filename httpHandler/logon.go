package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/logger"

	"github.com/rsmaxwell/players-api/session"
)

// Authenticate Response
type authenticateResponse struct {
	Token string `json:"token"`
}

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeLoginResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	logger.Logger.Printf("writeLoginResponse(0): user:%s, password:%s", user, pass)

	if !checkUser(user, pass) {
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

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(authenticateResponse{
		Token: token,
	})
}
