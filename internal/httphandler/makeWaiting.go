package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionMakeWaiting = debug.NewFunction(pkg, "MakeWaiting")
)

// MakeWaitingRequest structure
type MakeWaitingRequest struct {
	ID int `json:"id"`
}

// MakeWaiting method
func MakeWaiting(w http.ResponseWriter, r *http.Request) {
	f := functionMakeWaiting

	_, err := checkAuthenticated(r)
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

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("Unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	p := model.Person{ID: id}
	exists, err := p.PersonExists(db)
	if err != nil {
		message := fmt.Sprintf("Unexpected error checking person [%d] exists", id)
		f.DumpError(err, message)
		writeResponseError(w, r, codeerror.NewInternalServerError(message))
		return
	}
	if !exists {
		message := fmt.Sprintf("person [%d] not found", id)
		writeResponseError(w, r, codeerror.NewBadRequest(message))
		return
	}

	err = model.MakeWaiting(db, id)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
