package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/rsmaxwell/players-api/internal/common"
)

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Token  string                 `json:"token"`
	Person map[string]interface{} `json:"person"`
}

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(rw http.ResponseWriter, req *http.Request) {
	f := functionUpdatePerson

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r UpdatePersonRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	id := mux.Vars(req)["id"]

	f.DebugVerbose("ID:     %s", id)
	f.DebugVerbose("Token:  %s", r.Token)
	f.DebugVerbose("Person: %v", r.Person)

	err = model.UpdatePerson(r.Token, id, r.Person)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
