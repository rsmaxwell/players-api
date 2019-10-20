package httphandler

import (
	"log"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
)

// NotFound method
func NotFound(rw http.ResponseWriter, req *http.Request) {

	log.Printf("NotFound:")
	log.Printf("    Method:   %s", req.Method)
	log.Printf("    URL:   %s", req.URL)
	log.Printf("    Proto:   %s", req.Proto)
	log.Printf("    Host:   %s", req.Host)

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusNotFound, "Not Found")
	common.MetricsData.ClientError++
}
