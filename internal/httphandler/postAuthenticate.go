package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/gorilla/sessions"
	"github.com/rsmaxwell/players-api/internal/model"
)

// PostAuthenticateResponse structure
type PostAuthenticateResponse struct {
	Message string              `json:"message"`
	Person  model.LimitedPerson `json:"person"`
}

// Create the JWT key used to create the signature
var (
	key   = []byte("<SESSION_SECRET_KEY>")
	store = sessions.NewCookieStore(key)
)

// Authenticate method
func Authenticate(w http.ResponseWriter, r *http.Request) {
	f := functionAuthenticate

	if r.Method == http.MethodOptions {
		f.DebugVerbose("returning from 'Options' request")
		writeResponseMessage(w, r, http.StatusOK, "", "ok")
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err := fmt.Errorf("Unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	// *********************************************************************
	// * Authenticate the user
	// *********************************************************************
	userName, password, _ := r.BasicAuth()

	f.DebugVerbose("userName: %s", userName)
	f.DebugVerbose("password: %s", password)

	p, err := model.FindPersonByUserName(db, userName)
	if err != nil {
		writeResponseError(w, r, codeerror.NewUnauthorized("Not Authorised"))
		return
	}

	err = p.Authenticate(db, password)
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
		Path:     "players-api",
		MaxAge:   3600 * 6,
		HttpOnly: true,
	}

	sess.Values["userID"] = p.ID
	sess.Values["authenticated"] = true
	sess.Values["expiresAt"] = time.Now().Add(time.Hour * 6).Unix()

	err = sess.Save(r, w)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", PostAuthenticateResponse{
		Message: "ok",
		Person:  *p.ToLimited(),
	})
}
