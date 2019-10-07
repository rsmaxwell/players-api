package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"
)

// LogonResponse structure
type LogonResponse struct {
	Token string `json:"token"`
}

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {

	id, password, _ := req.BasicAuth()
	token, err := model.Login(id, password)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(LogonResponse{
		Token: token,
	})
}
