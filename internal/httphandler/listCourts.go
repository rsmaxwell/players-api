package httphandler

import (
	"database/sql"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionListCourts = debug.NewFunction(pkg, "ListCourts")
)

// ListCourts method
func ListCourts(w http.ResponseWriter, r *http.Request) {
	f := functionListCourts
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

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	listOfCourts, err := model.ListCourts(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, listOfCourts)
}
