package httphandler

import (
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

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	listOfCourts, err := model.ListCourts()
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", ListCourtsResponse{
		Courts: listOfCourts,
	})
}
