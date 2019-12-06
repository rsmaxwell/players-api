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
func GetCourt(w http.ResponseWriter, r *http.Request) {
	f := functionGetCourt

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	c, err := model.GetCourt(id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetCourtResponse{
		Court: *c,
	})
}
