package httphandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
)

// DeletePerson method
func DeletePerson(rw http.ResponseWriter, req *http.Request) {
	f := functionDeletePerson

	claims, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	if claims.UserID == id {
		f.DebugVerbose("Attempt delete self: %s", id)
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	err = model.DeletePerson(id)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
