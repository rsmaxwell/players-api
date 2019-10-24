package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// LogonResponse structure
type LogonResponse struct {
	Token string `json:"token"`
}

var (
	functionLogin = debug.NewFunction(pkg, "Login")
)

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {
	f := functionLogin

	id, password, _ := req.BasicAuth()

	f.DebugVerbose("id:       %s", id)
	f.DebugVerbose("password: %s", password)

	token, err := model.Login(id, password)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	f.DebugVerbose("token:    %s", token)

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(LogonResponse{
		Token: token,
	})
}
