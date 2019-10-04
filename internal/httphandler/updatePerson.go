package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Token  string                 `json:"token"`
	Person map[string]interface{} `json:"person"`
}

// UpdatePerson method
func UpdatePerson(rw http.ResponseWriter, req *http.Request, id string) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r UpdatePersonRequest
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

	err = model.UpdatePerson(id, session, r.Person)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
