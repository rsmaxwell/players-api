package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// ListCourtsResponse structure
type ListCourtsResponse struct {
	Courts []string `json:"courts"`
}

var (
	functionListCourts = debug.NewFunction(pkg, "GetQListCourtsueue")
)

// ListCourts method
func ListCourts(rw http.ResponseWriter, req *http.Request) {
	f := functionListCourts
	f.DebugVerbose("")

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

	listOfCourts, err := model.ListCourts()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(ListCourtsResponse{
		Courts: listOfCourts,
	})
}
