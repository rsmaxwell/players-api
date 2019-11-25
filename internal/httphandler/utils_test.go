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

func getLoginToken(t *testing.T, id, password string) error {

	// Create a  request to pass to our handler.
	req, err := http.NewRequest("GET", contextPath+"/login", nil)
	require.Nil(t, err, "err should be nothing")

	req.Header.Set("Authorization", model.BasicAuth(id, password))

	// Pass the request to our handler
	router := mux.NewRouter()
	SetupHandlers(router)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	require.Equal(t, http.StatusOK, rw.Code, "Error logging in: got %v want %v", http.StatusOK, rw.Code)

	return nil
}
