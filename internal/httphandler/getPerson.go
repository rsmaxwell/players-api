package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetPersonResponse structure
type GetPersonResponse struct {
	Person person.Person `json:"person"`
}

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(rw http.ResponseWriter, req *http.Request) {
	f := functionGetPerson

	session, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, err.Error())
		common.MetricsData.ServerError++
		return
	}
	defer session.SessionRelease(rw)
	userID := session.Get("id")
	if userID == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	p, err := model.GetPerson(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetPersonResponse{
		Person: *p,
	})
}
