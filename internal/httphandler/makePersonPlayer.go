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
func MakePersonPlayer(writer http.ResponseWriter, request *http.Request) {
	f := functionMakePersonPlayer

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	str := mux.Vars(request)["id1"]
	personID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}

	DebugVerbose(f, request, "PersonID: %d", personID)

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	err = model.MakePersonPlayerTx(db, personID)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
