package httphandler

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/rsmaxwell/players-api/internal/model"
)

// Create the JWT key used to create the signature
var (
	jwtKey = []byte("<JWT_SECRET_KEY>")
	key    = []byte("<SESSION_SECRET_KEY>")
	store  = sessions.NewCookieStore(key)
)

// Authenticate method
func Authenticate(rw http.ResponseWriter, req *http.Request) {
	f := functionAuthenticate

	id, password, _ := req.BasicAuth()

	f.DebugVerbose("id:       %s", id)
	f.DebugVerbose("password: %s", password)

	// *********************************************************************
	// * Authenticate the user
	// *********************************************************************
	p, err := model.Authenticate(id, password)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	p.Count++

	err = p.Save(id)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	// *********************************************************************
	// * Create a new session
	// *********************************************************************
	sess, err := store.Get(req, "players-api")
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	sess.Options = &sessions.Options{
		MaxAge:   3600 * 6,
		HttpOnly: true,
	}

	sess.Values["userID"] = id
	sess.Values["authenticated"] = true
	sess.Values["expiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	sess.Save(req, rw)

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
