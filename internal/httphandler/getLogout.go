package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionLogout = debug.NewFunction(pkg, "Logout")
)

// Logout method
func Logout(rw http.ResponseWriter, req *http.Request) {
	f := functionLogout
	f.DebugVerbose("")

	sess, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, err.Error())
		common.MetricsData.ServerError++
		return
	}
	defer sess.SessionRelease(rw)
	value := sess.Get("id")
	if value == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	globalSessions.SessionDestroy(rw, req)

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
