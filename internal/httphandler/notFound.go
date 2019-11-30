package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionNotFound = debug.NewFunction(pkg, "NotFound")
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {
	f := functionNotFound
	f.DebugVerbose("")

	writeResponseMessage(rw, req, http.StatusNotFound, "", "Not Found")
}
