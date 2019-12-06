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
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	f := functionGetMetrics

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	p, err := person.Load(userID)
	if err != nil {
		f.Dump("Could not load the logged on user: %s", userID)
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Not Authorized")
		return
	}

	if !p.CanGetMetrics() {
		f.DebugVerbose("unauthorized person[%s] attempted to get metrics", userID)
		writeResponseMessage(w, r, http.StatusForbidden, "", "Forbidden")
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetMetricsResponse{
		Data: common.MetricsData,
	})
}
