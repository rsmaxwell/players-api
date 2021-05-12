package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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

	str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		f.DebugVerbose("Count not convert '" + str + "' into an int")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	f.DebugVerbose("ID: %d", id)

	userID, ok := sess.Values["userID"].(int)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	if userID == id {
		f.DebugVerbose("Attempt delete self: %d", id)
		writeResponseMessage(w, r, http.StatusUnauthorized, "", "Not Authorized")
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	err = model.DeletePerson(db, id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
