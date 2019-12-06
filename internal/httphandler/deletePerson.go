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
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	f := functionDeletePerson

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	if userID == id {
		f.DebugVerbose("Attempt delete self: %s", id)
		writeResponseMessage(w, r, http.StatusUnauthorized, "", "Not Authorized")
		return
	}

	err = model.DeletePerson(id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
