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

// SigninRequest structure
type SigninResponse struct {
	Message     string       `json:"message"`
	Person      model.Person `json:"person"`
	AccessToken string       `json:"accessToken"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	functionBasicAuth = debug.NewFunction(pkg, "BasicAuth")
)

// BasicAuth function
func BasicAuth(username, password string) string {
	f := functionBasicAuth
	f.DebugVerbose("username: [%s], password:[%s]", username, password)

	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetSigninToken function
func GetSigninToken(t *testing.T, db *sql.DB, userName string, password string) (*http.Cookie, string) {

	command := "/signin"

	requestBody, err := json.Marshal(SigninRequest{
		Signin: model.Signin{
			Username: userName,
			Password: password,
		},
	})
	require.Nil(t, err, "err should be nothing")

	r, err := http.NewRequest("POST", contextPath+command, bytes.NewBuffer(requestBody))
	require.Nil(t, err, "err should be nothing")

	r.Header.Set("Authorization", BasicAuth(userName, password))
	w := httptest.NewRecorder()

	// ---------------------------------------

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
	defer cancel()
	r2 := r.WithContext(ctx)

	ctx = context.WithValue(r2.Context(), ContextDatabaseKey, db)
	r3 := r.WithContext(ctx)

	// ---------------------------------------

	router := mux.NewRouter()
	SetupHandlers(router)
	router.ServeHTTP(w, r3)

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
