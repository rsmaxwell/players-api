package httphandler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionRefresh = debug.NewFunction(pkg, "Refresh")
)

// Refresh method
func Refresh(rw http.ResponseWriter, req *http.Request) {
	f := functionRefresh

	// *********************************************************************
	// * Check the existing access token
	// *********************************************************************
	claims, err := checkAuthToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	// *********************************************************************
	// * Make a new access token and save it in a header
	// *********************************************************************
	accessToken := jwt.New(jwt.SigningMethodHS256)

	setAccessClaims(accessToken, claims)

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	f.DebugVerbose("Access-Token: %s", accessTokenString)
	rw.Header().Set("Access-Token", accessTokenString)

	// *********************************************************************
	// * Make a new refresh token and save it in a cookie
	// *********************************************************************
	refreshExpiration := time.Now().Add(time.Hour * 24)

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["sub"] = claims.UserID
	refreshClaims["exp"] = refreshExpiration.Unix()

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
	})

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
