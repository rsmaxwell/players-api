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
	f.DebugAPI("")

	userID, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}

	f.DebugVerbose("ID: %d", id)

	if userID == id {
		f.DebugVerbose("Attempt delete self: %d", id)
		writeResponseMessage(w, r, http.StatusUnauthorized, "Not Authorized")
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		writeResponseError(w, r, err)
		return
	}

	p := model.FullPerson{ID: id}
	err = p.DeletePerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
