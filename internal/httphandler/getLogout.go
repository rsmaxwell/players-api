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

	_, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "players-api",
		Value:    "",
		Expires:  time.Now(),
		MaxAge:   0,
		HttpOnly: true,
	})

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
