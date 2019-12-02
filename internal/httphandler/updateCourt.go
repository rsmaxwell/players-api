package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/person"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdateCourtRequest structure
type UpdateCourtRequest struct {
	Court map[string]interface{} `json:"court"`
}

var (
	functionUpdateCourt = debug.NewFunction(pkg, "UpdateCourt")
)

// UpdateCourt method
func UpdateCourt(rw http.ResponseWriter, req *http.Request) {
	f := functionUpdateCourt

	claims, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	p, err := person.Load(claims.UserID)
	if err != nil {
		f.Dump("Could not load the logged on user[%s]: %v", claims.UserID, err)
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", err.Error())
		return
	}
	if !p.CanUpdateCourt() {
		writeResponseMessage(rw, req, http.StatusUnauthorized, "", "Not Authorized")
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var r UpdateCourtRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdateCourt(id, r.Court)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
