package httphandler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/stretchr/testify/require"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	functionBasicAuth = debug.NewFunction(pkg, "BasicAuth")
)

// BasicAuth function
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetSigninToken function
func GetSigninToken(t *testing.T, db *sql.DB, userName string, password string) (*http.Cookie, string) {
	f := functionBasicAuth

	command := "/signin"

	requestBody, err := json.Marshal(SigninRequest{
		Signin: model.Signin{
			Username: userName,
			Password: password,
		},
	})
	require.Nil(t, err, "err should be nothing")

	request, err := http.NewRequest("POST", contextPath+command, bytes.NewBuffer(requestBody))
	require.Nil(t, err, "err should be nothing")

	DebugVerbose(f, request, "username: [%s], password:[%s]", userName, password)

	request.Header.Set("Authorization", BasicAuth(userName, password))
	w := httptest.NewRecorder()

	// ---------------------------------------

	ctx, cancel := context.WithTimeout(request.Context(), time.Duration(60*time.Second))
	defer cancel()
	request2 := request.WithContext(ctx)

	ctx = context.WithValue(request2.Context(), ContextDatabaseKey, db)
	request3 := request.WithContext(ctx)

	// ---------------------------------------

	router := mux.NewRouter()
	SetupHandlers(router)
	router.ServeHTTP(w, request3)

	response := w.Result()
	require.Equal(t, http.StatusOK, w.Code, "authentication failed")

	cookies := map[string]*http.Cookie{}
	for _, cookie := range response.Cookies() {
		cookies[cookie.Name] = cookie
	}

	cookie := cookies["players-api"]
	require.NotNil(t, cookie, "cookie missing")

	bytes, err := ioutil.ReadAll(w.Body)
	require.Nil(t, err, "err should be nothing")

	responseBody := new(SigninResponse)
	err = json.Unmarshal(bytes, &responseBody)
	require.Nil(t, err, "err should be nothing")

	token := responseBody.AccessToken
	require.NotNil(t, token, "token missing")

	return cookie, token
}
