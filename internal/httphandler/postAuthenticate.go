package httphandler

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/rsmaxwell/players-api/internal/model"
)

// SigninRequest structure
type SigninRequest struct {
	Signin model.Signin `json:"signin"`
}

// PostSigninResponse structure
type PostSigninResponse struct {
	Message     string       `json:"message"`
	Person      model.Person `json:"person"`
	AccessToken string       `json:"accessToken"`
}

// Signin method
func Signin(w http.ResponseWriter, r *http.Request) {
	f := functionSignin
	f.DebugAPI("")

	if r.Method == http.MethodOptions {
		writeResponseMessage(w, r, http.StatusOK, "ok")
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request SigninRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// *********************************************************************
	// * Authenticate the user
	// *********************************************************************
	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	email := request.Signin.Username

	p, err := model.FindPersonByEmail(db, email)
	if err != nil {
		writeResponseError(w, r, codeerror.NewUnauthorized("Not Authenticated"))
		return
	}

	err = p.Authenticate(db, request.Signin.Password)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	// *********************************************************************
	// * Create the token pair
	// *********************************************************************
	newAccessToken, err := basic.GenerateToken(p.ID, time.Minute*time.Duration(5))
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	newRefreshToken, err := basic.GenerateToken(p.ID, time.Hour*time.Duration(1))
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
		MaxAge:   3600 * 1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["userID"] = p.ID
	sess.Values["refreshToken"] = newRefreshToken
	sess.Values["expiresAt"] = time.Now().Add(time.Hour * 1).Unix()

	err = sess.Save(r, w)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	// *********************************************************************
	// * Write the response
	// *********************************************************************
	limitedPerson := p.ToLimited()
	writeResponseObject(w, r, http.StatusOK, PostSigninResponse{
		Message:     "ok",
		Person:      *limitedPerson,
		AccessToken: newAccessToken,
	})
}
