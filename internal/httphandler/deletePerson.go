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

	sess, err := checkAuthenticated(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "Error")
		return
	}

	if userID == id {
		f.DebugVerbose("Attempt delete self: %s", id)
		writeResponseMessage(rw, req, http.StatusUnauthorized, "", "Not Authorized")
		return
	}

	err = model.DeletePerson(id)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
