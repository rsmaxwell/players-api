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
	functionClearCourt = debug.NewFunction(pkg, "ClearCourt")
)

// ClearCourt method
func ClearCourt(w http.ResponseWriter, r *http.Request) {
	f := functionClearCourt
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

	err = model.ClearCourt(db, courtID)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}
	writeResponseObject(w, r, http.StatusOK, response)
}
