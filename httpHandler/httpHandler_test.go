package httpHandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assert that the response is JSON of the expected format
func assertResponseJSON(t *testing.T, rr *httptest.ResponseRecorder, expectedObj interface{}) {
	rr.Flush()
	assert.Equal(t, rr.Header()["Content-Type"], []string{"application/json"}, "Unexpected Content-Type")
	ej, err := json.Marshal(expectedObj)
	assert.Nil(t, err)
	assert.JSONEq(t, string(ej), string(rr.Body.Bytes()))
}

// Testing the writeStatusResponse Function
func TestWriteMessageResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	msg := "A message of sorts"
	status := http.StatusInternalServerError

	expectedObj := messageResponse{
		Message: msg,
	}

	WriteResponse(rr, status, msg)
	assert.Equal(t, status, rr.Code, "Wrong HTTP status code returned")
	assertResponseJSON(t, rr, expectedObj)
}
