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
func ClearCourt(writer http.ResponseWriter, request *http.Request) {
	f := functionClearCourt

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	str := mux.Vars(request)["id"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}
	DebugVerbose(f, request, "courtID: %d", courtID)

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	err = model.ClearCourtTx(db, courtID)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}
	writeResponseObject(writer, request, http.StatusOK, response)
}
