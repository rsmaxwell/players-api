package httphandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionDeleteCourt = debug.NewFunction(pkg, "DeleteCourt")
)

// DeleteCourt method
func DeleteCourt(rw http.ResponseWriter, req *http.Request) {
	f := functionDeleteCourt

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.DeleteCourt(id)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
