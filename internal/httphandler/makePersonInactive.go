package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionMakePersonInactive = debug.NewFunction(pkg, "MakePersonInactive")
)

// MakeInactiveRequest structure
type MakeInactiveRequest struct {
	ID int `json:"id"`
}

// MakePersonInactive method
func MakePersonInactive(w http.ResponseWriter, r *http.Request) {
	f := functionMakePersonInactive
	f.DebugAPI("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id"]
	personID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	f.DebugVerbose("ID: %d", personID)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	err = model.MakePersonInactive(db, personID)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
