package httphandler

import (
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

	sess, err := checkAuthenticated(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "Error")
		return
	}

	p, err := person.Load(userID)
	if err != nil {
		f.Dump("Could not load the logged on user: %s", userID)
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "Not Authorized")
		return
	}

	if !p.CanGetMetrics() {
		f.DebugVerbose("unauthorized person[%s] attempted to get metrics", userID)
		writeResponseMessage(rw, req, http.StatusForbidden, "", "Forbidden")
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", GetMetricsResponse{
		Data: common.MetricsData,
	})
}
