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
func Register(rw http.ResponseWriter, req *http.Request) {
	f := functionRegister

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var r RegisterRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugVerbose("UserID:    %s", r.UserID)
	f.DebugVerbose("Password:  %s", r.Password)
	f.DebugVerbose("FirstName: %s", r.FirstName)
	f.DebugVerbose("LastName:  %s", r.LastName)
	f.DebugVerbose("Email:     %s", r.Email)

	err = model.Register(r.UserID, r.Password, r.FirstName, r.LastName, r.Email)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
