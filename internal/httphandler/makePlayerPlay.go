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
func MakePlayerPlay(writer http.ResponseWriter, request *http.Request) {
	f := functionMakePlayerPlay

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

	str = mux.Vars(request)["id2"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	DebugVerbose(f, request, "courtID: %d", courtID)

	str = mux.Vars(request)["id3"]
	position, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	DebugVerbose(f, request, "position:%d", position)

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	if position < 0 {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("unexpected position: [%d]", position))
		return
	}

	if position >= model.NumberOfCourtPositions {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("unexpected position: [%d]", position))
		return
	}

	err = model.MakePlayerPlayTx(db, personID, courtID, position)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
