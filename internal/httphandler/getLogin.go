package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
)

// LogonResponse structure
type LogonResponse struct {
	Token string `json:"token"`
}

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {

	// Check the userID calling the service
	userID, pass, _ := req.BasicAuth()
	if !model.CheckUser(userID, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	if !model.PersonCanLogin(userID) {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	token, err := session.New(userID)
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
