package httphandler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionRefresh = debug.NewFunction(pkg, "Refresh")
)

// Refresh method
func Refresh(rw http.ResponseWriter, req *http.Request) {
	f := functionRefresh

	// *********************************************************************
	// * Check the existing refresh token and count
	// *********************************************************************
	refreshClaims, err := checkRefreshToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	p, err := person.Load(refreshClaims.UserID)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	if p.Count != refreshClaims.Count {
		if err != nil {
			writeResponseError(rw, req, err)
			return
		}
	}

	// *********************************************************************
	// * Make a new access token and save it in a header
	// *********************************************************************
	accessToken := jwt.New(jwt.SigningMethodHS256)

	setAccessClaims(accessToken, &AccessClaims{
		UserID:    refreshClaims.UserID,
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
	// * Make a new refresh token and save it in a cookie
	// *********************************************************************
	refreshExpiration := time.Now().Add(time.Hour * 24)

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	refreshClaims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	refreshClaims.Count++
	setRefreshClaims(refreshToken, refreshClaims)

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

	// *********************************************************************
	// * Update the count
	// *********************************************************************
	p.Count++
	err = p.Save(refreshClaims.UserID)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
