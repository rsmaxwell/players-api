package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {
	f := debug.NewFunction(pkg, "NotFound")
	f.DebugVerbose("")

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusNotFound, "Not Found")
	common.MetricsData.ClientError++
}
