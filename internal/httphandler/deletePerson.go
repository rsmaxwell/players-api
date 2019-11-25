package httphandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
)

// DeletePerson method
func DeletePerson(rw http.ResponseWriter, req *http.Request) {
	f := functionDeletePerson

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
		f.Dump("Unexpected type for userID: %T: %v", value, value)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	if userID == id {
		f.DebugVerbose("Attempt delete self: %s", id)
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	err = model.DeletePerson(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
