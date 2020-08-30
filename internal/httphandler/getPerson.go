package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetPersonResponse structure
type GetPersonResponse struct {
	Message string              `json:"message"`
	Person  model.LimitedPerson `json:"person"`
}

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(w http.ResponseWriter, r *http.Request) {
	f := functionGetPerson

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		message := "Could not convert '" + str + "' into an int"
		f.Errorf(message)
		err = codeerror.NewBadRequest(message)
		writeResponseError(w, r, err)
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

	var p model.Person
	p.ID = id
	err = p.LoadPerson(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetPersonResponse{
		Message: "ok",
		Person:  *p.ToLimited(),
	})
}
