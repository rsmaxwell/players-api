package httphandler

import (
	"net/http"
	"time"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionLogout = debug.NewFunction(pkg, "Logout")
)

// Logout method
func Logout(rw http.ResponseWriter, req *http.Request) {
	f := functionLogout
	f.DebugVerbose("")

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "players-api",
		Value:    "",
		Expires:  time.Now(),
		MaxAge:   0,
		HttpOnly: true,
	})

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
