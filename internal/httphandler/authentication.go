package httphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"

	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionSignin             = debug.NewFunction(pkg, "Signin")
	functionCheckAuthenticated = debug.NewFunction(pkg, "checkAuthenticated")
)

// SigninRequest structure
type SigninRequest struct {
	Signin model.Signin `json:"signin"`
}

// SigninResponse structure
type SigninResponse struct {
	Message      string       `json:"message"`
	Person       model.Person `json:"person"`
	AccessToken  string       `json:"accessToken"`
	RefreshDelta int          `json:"refreshDelta"`
}

// Signin method
func Signin(writer http.ResponseWriter, request *http.Request) {
	f := functionSignin

	limitedReader := &io.LimitedReader{R: request.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	DebugRequestBody(f, request, b)

	var signinRequest SigninRequest
	err = json.Unmarshal(b, &signinRequest)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	// *********************************************************************
	// * Authenticate the user
	// *********************************************************************
	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := fmt.Sprintf("unexpected context type: %#v", db)
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	email := signinRequest.Signin.Username

	p, err := model.FindPersonByEmail(context.Background(), db, email)
	if err != nil {
		writeResponseError(writer, request, codeerror.NewUnauthorized("Not Authenticated"))
		return
	}

	err = p.Authenticate(db, signinRequest.Signin.Password)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	// *********************************************************************
	// * Create the token pair
	// *********************************************************************
	object = request.Context().Value(ContextConfigKey)
	cfg, ok := object.(*config.Config)
	if !ok {
		message := fmt.Sprintf("unexpected context type: %#v", cfg)
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	f.DebugVerbose("accessTokenExpiry:  %10s     expires at: %s", cfg.AccessTokenExpiry, time.Now().Add(cfg.AccessTokenExpiry).Round(time.Second))
	newAccessToken, err := basic.GenerateToken(p.ID, cfg.AccessTokenExpiry)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	f.DebugVerbose("refreshTokenExpiry: %10s     expires at: %s", cfg.RefreshTokenExpiry, time.Now().Add(cfg.RefreshTokenExpiry).Round(time.Second))
	newRefreshToken, err := basic.GenerateToken(p.ID, cfg.RefreshTokenExpiry)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	f.DebugVerbose("clientRefreshDelta: %10s", cfg.ClientRefreshDelta)

	// *********************************************************************
	// * Create a new session to contain the refresh token
	// *********************************************************************
	sess, err := store.Get(request, "players-api")
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(cfg.RefreshTokenExpiry / time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["userID"] = p.ID
	sess.Values["refreshToken"] = newRefreshToken
	sess.Values["expiresAt"] = time.Now().UTC().Add(cfg.RefreshTokenExpiry).Unix()

	err = sess.Save(request, writer)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	// *********************************************************************
	// * Write the response
	// *********************************************************************
	limitedPerson := p.ToLimited()
	writeResponseObject(writer, request, http.StatusOK, SigninResponse{
		Message:      "ok",
		Person:       *limitedPerson,
		AccessToken:  newAccessToken,
		RefreshDelta: int(cfg.ClientRefreshDelta / time.Second),
	})
}

func NewDuration(hours int, minutes int, seconds int) (time.Duration, int) {
	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	durationInSeconds := hours*3600 + minutes*60 + seconds
	return duration, durationInSeconds
}

// Signout method
func Signout(w http.ResponseWriter, req *http.Request) {

	_, err := checkAuthenticated(req)
	if err != nil {
		writeResponseError(w, req, err)
		return
	}

	writeResponseMessage(w, req, http.StatusOK, "ok")
}

// checkAuthenticated method
func checkAuthenticated(request *http.Request) (int, error) {
	f := functionCheckAuthenticated

	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		DebugError(f, request, "missing Authorization header")
		return 0, fmt.Errorf("not authorized")
	}

	splitToken := strings.Split(authorizationHeader, "Bearer ")
	if splitToken == nil {
		return 0, fmt.Errorf("not authorized")
	}

	tokenString := splitToken[1]

	claims, err := basic.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	DebugVerbose(f, request, fmt.Sprintf("user: %d", claims.ID))

	return claims.ID, nil
}
