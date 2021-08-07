package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionNotFound = debug.NewFunction(pkg, "NotFound")
)

// NotFound method
func NotFound(w http.ResponseWriter, r *http.Request) {
	f := functionNotFound
	f.DebugAPI("")

	writeResponseMessage(w, r, http.StatusNotFound, "Not Found")
}
