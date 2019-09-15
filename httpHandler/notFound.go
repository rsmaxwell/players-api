package httpHandler

import (
	"net/http"
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {
	WriteResponse(rw, http.StatusNotFound, "Not Found")
	clientError++
}
