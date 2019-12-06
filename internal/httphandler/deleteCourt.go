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
func DeleteCourt(w http.ResponseWriter, r *http.Request) {
	f := functionDeleteCourt

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.DeleteCourt(id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
