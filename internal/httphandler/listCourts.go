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
func ListCourts(w http.ResponseWriter, r *http.Request) {
	f := functionListCourts
	f.DebugVerbose("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	listOfCourts, err := model.ListCourts()
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", ListCourtsResponse{
		Courts: listOfCourts,
	})
}
