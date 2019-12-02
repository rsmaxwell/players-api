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

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Person map[string]interface{} `json:"person"`
}

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(rw http.ResponseWriter, req *http.Request) {
	f := functionUpdatePerson

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

	var r UpdatePersonRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(req)["id"]

	f.DebugVerbose("ID:     %s", id)
	f.DebugVerbose("Person: %v", r.Person)

	err = model.UpdatePerson(claims.UserID, id, r.Person)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
