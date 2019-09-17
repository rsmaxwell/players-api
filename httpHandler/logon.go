package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/session"
)

// AuthenticateResponse structure
type AuthenticateResponse struct {
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
	json.NewEncoder(rw).Encode(AuthenticateResponse{
		Token: token,
	})
}
