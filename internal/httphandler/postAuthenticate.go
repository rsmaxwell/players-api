package httphandler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/rsmaxwell/players-api/internal/model"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("<JWT_SECRET_KEY>")

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
	// * Make an access token and put it in the header
	// *********************************************************************
	accessToken := jwt.New(jwt.SigningMethodHS256)

	setAccessClaims(accessToken, &AccessClaims{
		UserID:    id,
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Role:      p.Role,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	f.DebugVerbose("Access-Token: %s", accessTokenString)
	rw.Header().Set("Access-Token", accessTokenString)

	// *********************************************************************
	// * Make a refresh token and put it in a cookie
	// *********************************************************************
	refreshExpiration := time.Now().Add(time.Hour * 24)

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	setRefreshClaims(refreshToken, &RefreshClaims{
		UserID:    id,
		ExpiresAt: refreshExpiration.Unix(),
		Count:     p.Count,
	})

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	f.DebugVerbose("Refresh-Token: %s", refreshTokenString)

	http.SetCookie(rw, &http.Cookie{
		Name:     "players-api",
		Value:    refreshTokenString,
		Expires:  refreshExpiration,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
