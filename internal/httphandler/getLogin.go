package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionLogin = debug.NewFunction(pkg, "Login")
)

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {
	f := functionLogin

	id, password, _ := req.BasicAuth()

	f.DebugVerbose("id:       %s", id)
	f.DebugVerbose("password: %s", password)

	err := model.Authenticate(id, password)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	sess, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}
	defer sess.SessionRelease(rw)

	sess.Set("id", id)

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
