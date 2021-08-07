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

// GetPersonResponse structure
type GetPersonResponse struct {
	Message string       `json:"message"`
	Person  model.Person `json:"person"`
}

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(w http.ResponseWriter, r *http.Request) {
	f := functionGetPerson
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

	f.DebugVerbose("ID: %d", id)

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	var p model.FullPerson
	p.ID = id
	err = p.LoadPerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, GetPersonResponse{
		Message: "ok",
		Person:  *p.ToLimited(),
	})
}
