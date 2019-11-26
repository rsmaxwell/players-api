package httphandler

import (
	"encoding/json"
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

	_, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	c, err := model.GetCourt(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetCourtResponse{
		Court: *c,
	})
}
