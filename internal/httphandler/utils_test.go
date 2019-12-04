package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var (
	goodUserID    = "007"
	goodPassword  = "topsecret"
	goodCourtID   = "1000"
	anotherUserID = "bob"
)

func TestGetLoginToken(t *testing.T) {

	teardown := model.SetupFull(t)
	defer teardown(t)

	testLogin(t, goodUserID, goodPassword)
}

func testLogin(t *testing.T, userID, password string) (string, *http.Cookie) {
	req, err := http.NewRequest("POST", contextPath+"/users/authenticate", nil)
	require.Nil(t, err, "err should be nothing")

	req.Header.Set("Authorization", model.BasicAuth(userID, password))

	router := mux.NewRouter()
	SetupHandlers(router)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	response := rw.Result()
	require.Equal(t, rw.Code, http.StatusOK, "authentication failed")

	cookies := map[string]*http.Cookie{}
	for _, cookie := range response.Cookies() {
		cookies[cookie.Name] = cookie
	}

	cookie := cookies["players-api"]
	require.NotNil(t, cookie, "cookie should be something")

	return "", cookie
}

func setAccessToken(req *http.Request, setAccessToken bool, accessToken string) {
	if setAccessToken {
		req.Header.Set("Authorization", accessToken)
	}
}

func setRefreshToken(req *http.Request, useGoodRefreshToken, setRefreshToken bool, refreshTokenCookie *http.Cookie, refreshToken string) {
	if useGoodRefreshToken {
		req.AddCookie(refreshTokenCookie)
	} else if setRefreshToken {

		refreshExpiration := time.Now().Add(time.Hour * 24)

		req.AddCookie(&http.Cookie{
			Name:     "players-api",
			Value:    refreshToken,
			Expires:  refreshExpiration,
			HttpOnly: true,
		})
	}
}
