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
	functionMakePlayerWait = debug.NewFunction(pkg, "MakePlayerWait")
)

// MakeWaitingRequest structure
type MakeWaitingRequest struct {
	ID int `json:"id"`
}

// MakeWaiting method
func MakePlayerWait(writer http.ResponseWriter, request *http.Request) {
	f := functionMakePlayerWait

	_, err := checkAuthenticated(request)
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

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	err = model.MakePlayerWaitTx(db, id)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
