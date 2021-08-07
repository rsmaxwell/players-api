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
	functionFillCourt = debug.NewFunction(pkg, "FillCourt")
)

// FillCourtRequest structure
type FillCourtRequest struct {
	PersonID int `json:"id"`
}

// FillCourt method
func FillCourt(w http.ResponseWriter, r *http.Request) {
	f := functionFillCourt
	f.DebugAPI("")

	if r.Method == http.MethodOptions {
		writeResponseMessage(w, r, http.StatusOK, "ok")
		return
	}

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	f.DebugVerbose("courtID: %d", courtID)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	positions, err := model.FillCourt(db, courtID)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	response := struct {
		Message   string           `json:"message"`
		Positions []model.Position `json:"positions"`
	}{
		Message:   "ok",
		Positions: positions,
	}
	f.DebugInfo(fmt.Sprintf("Response: %v", response))
	writeResponseObject(w, r, http.StatusOK, response)
}
