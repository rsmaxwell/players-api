package httphandler

import (
	"net/http"
)

// NotFound method
func NotFound(w http.ResponseWriter, r *http.Request) {
	writeResponseMessage(w, r, http.StatusNotFound, "Not Found")
}
