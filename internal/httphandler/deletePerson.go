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

var (
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
)

// DeletePerson method
func DeletePerson(writer http.ResponseWriter, request *http.Request) {
	f := functionDeletePerson

	userID, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	str := mux.Vars(request)["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}

	DebugVerbose(f, request, "ID: %d", id)

	if userID == id {
		DebugVerbose(f, request, "Attempt delete self: %d", id)
		writeResponseMessage(writer, request, http.StatusUnauthorized, "Not Authorized")
		return
	}

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		writeResponseError(writer, request, err)
		return
	}

	p := model.FullPerson{ID: id}
	err = p.DeletePersonTx(db)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
