package httphandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetCourtResponse structure
type GetCourtResponse struct {
	Court court.Court `json:"court"`
}

var (
	functionGetCourt = debug.NewFunction(pkg, "GetCourt")
)

// GetCourt method
func GetCourt(rw http.ResponseWriter, req *http.Request) {
	f := functionGetCourt

	_, err := checkAuthenticated(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	c, err := model.GetCourt(id)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", GetCourtResponse{
		Court: *c,
	})
}
