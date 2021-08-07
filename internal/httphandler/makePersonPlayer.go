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
	functionMakePersonPlayer = debug.NewFunction(pkg, "MakePersonPlayer")
)

// MakePlayingRequest structure
type MakePlayingRequest struct {
	PersonID int `json:"id"`
}

// MakePersonPlayer method
func MakePersonPlayer(w http.ResponseWriter, r *http.Request) {
	f := functionMakePersonPlayer
	f.DebugAPI("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id1"]
	personID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	f.DebugVerbose("PersonID: %d", personID)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	err = model.MakePersonPlayer(db, personID)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
