package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdatePersonRoleRequest structure
type UpdatePersonRoleRequest struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// UpdatePersonRole method
func UpdatePersonRole(rw http.ResponseWriter, req *http.Request, id string) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r UpdatePersonRoleRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	if !model.PersonCanUpdatePersonRole(session.UserID, id) {
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		clientError++
		return
	}

	err = model.UpdatePersonRole(id, r.Role)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
