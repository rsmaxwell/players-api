package httphandler

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rsmaxwell/players-api/internal/model"
)

// Create the JWT key used to create the signature
var (
	key   = []byte("<SESSION_SECRET_KEY>")
	store = sessions.NewCookieStore(key)
)

// Authenticate method
func Authenticate(w http.ResponseWriter, r *http.Request) {
	f := functionAuthenticate

	id, password, _ := r.BasicAuth()

	f.DebugVerbose("id:       %s", id)
	f.DebugVerbose("password: %s", password)

	// *********************************************************************
	// * Authenticate the user
	// *********************************************************************
	_, err := model.Authenticate(id, password)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	// *********************************************************************
	// * Create a new session
	// *********************************************************************
	sess, err := store.New(r, "players-api")
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	sess.Options = &sessions.Options{
		MaxAge:   3600 * 6,
		HttpOnly: true,
	}

	sess.Values["userID"] = id
	sess.Values["authenticated"] = true
	sess.Values["expiresAt"] = time.Now().Add(time.Hour * 24).Unix()

	err = sess.Save(r, w)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
