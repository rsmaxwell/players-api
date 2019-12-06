package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionLogout = debug.NewFunction(pkg, "Logout")
)

// Logout method
func Logout(w http.ResponseWriter, r *http.Request) {
	f := functionLogout
	f.DebugVerbose("")

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	sess.Options.MaxAge = -1

	err = sess.Save(r, w)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
