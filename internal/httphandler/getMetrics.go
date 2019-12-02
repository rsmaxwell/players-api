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

	claims, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	p, err := person.Load(claims.UserID)
	if err != nil {
		f.Dump("Could not load the logged on user: %s", claims.UserID)
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "Not Authorized")
		return
	}

	if !p.CanGetMetrics() {
		f.DebugVerbose("unauthorized person[%s] attempted to get metrics", claims.UserID)
		writeResponseMessage(rw, req, http.StatusForbidden, "", "Forbidden")
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", GetMetricsResponse{
		Data: common.MetricsData,
	})
}
