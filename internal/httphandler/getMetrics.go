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

	claims, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	p, err := person.Load(claims.UserID)
	if err != nil {
		f.Dump("Could not load the logged on user: %s", claims.UserID)
		WriteResponse(rw, http.StatusInternalServerError, "Not Authorized")
		return
	}

	if !p.CanGetMetrics() {
		f.DebugVerbose("unauthorized person[%s] attempted to get metrics", claims.UserID)
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetMetricsResponse{
		Data: common.MetricsData,
	})
}
