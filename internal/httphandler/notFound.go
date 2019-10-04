package httphandler

import (
	"net/http"
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {
	setHeaders(rw, req)
	WriteResponse(rw, http.StatusNotFound, "Not Found")
	clientError++
}
