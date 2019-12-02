package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdatePersonRoleRequest structure
type UpdatePersonRoleRequest struct {
	Role string `json:"role"`
}

var (
	functionUpdatePersonRole = debug.NewFunction(pkg, "UpdatePersonRole")
)

// UpdatePersonRole method
func UpdatePersonRole(rw http.ResponseWriter, req *http.Request) {
	f := functionUpdatePersonRole

	claims, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var r UpdatePersonRoleRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdatePersonRole(claims.UserID, id, r.Role)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
