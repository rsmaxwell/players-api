package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func testLogin(t *testing.T, userID, password string) *http.Cookie {
	r, err := http.NewRequest("POST", contextPath+"/users/authenticate", nil)
	require.Nil(t, err, "err should be nothing")

	r.Header.Set("Authorization", model.BasicAuth(userID, password))
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	SetupHandlers(router)
	router.ServeHTTP(w, r)

	response := w.Result()
	require.Equal(t, w.Code, http.StatusOK, "authentication failed")

	cookies := map[string]*http.Cookie{}
	for _, cookie := range response.Cookies() {
		cookies[cookie.Name] = cookie
	}

	cookie := cookies["players-api"]
	return cookie
}
