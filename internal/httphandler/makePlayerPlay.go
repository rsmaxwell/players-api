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
	functionMakePlayerPlay = debug.NewFunction(pkg, "MakePlayerPlay")
)

// MakePlayerPlayRequest structure
type MakePlayerPlayRequest struct {
	PersonID int `json:"id"`
}

// MakePlayerPlay method
func MakePlayerPlay(w http.ResponseWriter, r *http.Request) {
	f := functionMakePlayerPlay
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

	str = mux.Vars(r)["id2"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	f.DebugVerbose("courtID: %d", courtID)

	str = mux.Vars(r)["id3"]
	position, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	f.DebugVerbose("position:%d", position)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	if position < 0 {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected position: [%d]", position))
		return
	}

	if position >= model.NumberOfCourtPositions {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("unexpected position: [%d]", position))
		return
	}

	err = model.MakePlayerPlay(db, personID, courtID, position)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
