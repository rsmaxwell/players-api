package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// GetMetricsResponse structure
type GetMetricsResponse struct {
	Data common.Metrics `json:"data"`
}

var (
	functionGetMetrics = debug.NewFunction(pkg, "GetMetrics")
)

// GetMetrics method
func GetMetrics(rw http.ResponseWriter, req *http.Request) {
	f := functionGetMetrics

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
		f.Dump("Unexpected type for userID: %t, %v", value, value)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}

	p, err := person.Load(userID)
	if err != nil {
		f.Dump("Could not load the logged on user: %s", userID)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}

	if !p.CanGetMetrics() {
		f.DebugVerbose("unauthorized person[%s] attempted to get metrics", userID)
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetMetricsResponse{
		Data: common.MetricsData,
	})
}
