package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// ListCourtsResponse structure
type ListCourtsResponse struct {
	Courts []string `json:"courts"`
}

var (
	functionListCourts = debug.NewFunction(pkg, "GetQListCourtsueue")
)

// ListCourts method
func ListCourts(rw http.ResponseWriter, req *http.Request) {
	f := functionListCourts
	f.DebugVerbose("")

	_, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	listOfCourts, err := model.ListCourts()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(ListCourtsResponse{
		Courts: listOfCourts,
	})
}
