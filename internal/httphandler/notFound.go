package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {
	setHeaders(rw, req)
	WriteResponse(rw, http.StatusNotFound, "Not Found")
	common.MetricsData.ClientError++
}
