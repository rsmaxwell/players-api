package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// RegisterRequest structure
type RegisterRequest struct {
	UserID    string `json:"userID"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

var (
	functionRegister = debug.NewFunction(pkg, "Register")
)

// Register method
func Register(w http.ResponseWriter, r *http.Request) {
	f := functionRegister

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request RegisterRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugVerbose("UserID:    %s", request.UserID)
	f.DebugVerbose("Password:  %s", request.Password)
	f.DebugVerbose("FirstName: %s", request.FirstName)
	f.DebugVerbose("LastName:  %s", request.LastName)
	f.DebugVerbose("Email:     %s", request.Email)

	err = model.Register(request.UserID, request.Password, request.FirstName, request.LastName, request.Email)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
