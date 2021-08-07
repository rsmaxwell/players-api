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

// GetCourtResponse structure
type GetCourtResponse struct {
	Message string      `json:"message"`
	Court   model.Court `json:"court"`
}

var (
	functionGetCourt = debug.NewFunction(pkg, "GetCourt")
)

// GetCourt method
func GetCourt(w http.ResponseWriter, r *http.Request) {
	f := functionGetCourt
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
	id, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
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

	var c model.Court
	c.ID = id
	err = c.LoadCourt(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, GetCourtResponse{
		Message: "ok",
		Court:   c,
	})
}
