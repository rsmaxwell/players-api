package httphandler

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/debug"
)

// GetRefreshTokensResponse structure
type GetRefreshTokensResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}

var (
	functionRefreshTokens = debug.NewFunction(pkg, "functionRefreshTokens")
)

// RefreshTokens method
func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	f := functionRefreshTokens
	f.DebugAPI("")

	session, err := store.Get(r, "players-api")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenObject := session.Values["refreshToken"]
	tokenString, ok := tokenObject.(string)
	if !ok {
		http.Error(w, "refreshToken not found", http.StatusInternalServerError)
		return
	}

	claims, err := basic.ValidateToken(tokenString)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	// *********************************************************************
	// * Create the token pair
	// *********************************************************************
	newAccessToken, err := basic.GenerateToken(claims.ID, time.Second*time.Duration(5))
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	newRefreshToken, err := basic.GenerateToken(claims.ID, time.Hour*time.Duration(1))
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
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["userID"] = claims.ID
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
	writeResponseObject(w, r, http.StatusOK, GetRefreshTokensResponse{
		Message:     "ok",
		AccessToken: newAccessToken,
	})

}
