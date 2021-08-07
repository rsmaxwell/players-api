package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionLogout = debug.NewFunction(pkg, "Signout")
)

// Signout method
func Signout(w http.ResponseWriter, r *http.Request) {
	f := functionLogout
	f.DebugAPI("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
