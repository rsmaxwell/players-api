package httphandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionRefreshToken = debug.NewFunction(pkg, "RefreshToken")
)

// GetRefreshTokensResponse structure
type GetRefreshTokensResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}

// RefreshToken method
func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	f := functionRefreshToken

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	session, err := store.Get(request, "players-api")
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	tokenObject := session.Values["refreshToken"]
	tokenString, ok := tokenObject.(string)
	if !ok {
		DebugVerbose(f, request, "refreshToken not found")
		err = codeerror.NewUnauthorized("Unauthorized")
		writeResponseError(writer, request, err)
		return
	}

	claims, err := basic.ValidateToken(tokenString)
	if err != nil {
		DebugVerbose(f, request, "refreshToken not valid: %s", err.Error())
		err = codeerror.NewUnauthorized("Unauthorized")
		writeResponseError(writer, request, err)
		return
	}

	// *********************************************************************
	// * Create a new access token
	// *********************************************************************
	object := request.Context().Value(ContextConfigKey)
	cfg, ok := object.(*config.Config)
	if !ok {
		message := fmt.Sprintf("unexpected context type: %#v", cfg)
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	f.DebugVerbose("accessTokenExpiry:  %10s     expires at: %s", cfg.AccessTokenExpiry, time.Now().Add(cfg.AccessTokenExpiry))
	newAccessToken, err := basic.GenerateToken(claims.ID, cfg.AccessTokenExpiry)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	// *********************************************************************
	// * Write the response
	// *********************************************************************
	writeResponseObject(writer, request, http.StatusOK, GetRefreshTokensResponse{
		Message:     "ok",
		AccessToken: newAccessToken,
	})

}
