package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/person"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/common"
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

	session, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, err.Error())
		common.MetricsData.ServerError++
		return
	}
	defer session.SessionRelease(rw)
	value := session.Get("id")
	if value == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}
	userID, ok := value.(string)
	if !ok {
		f.Dump("Unexpected type for userID: %t: %v", value, value)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}
	p, err := person.Load(userID)
	if err != nil {
		f.Dump("Could not load the logged on user[%s]: %v", userID, err)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}
	if !p.CanUpdateCourt() {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r UpdateCourtRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdateCourt(id, r.Court)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
