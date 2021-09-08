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
	functionFillCourt = debug.NewFunction(pkg, "FillCourt")
)

// FillCourtRequest structure
type FillCourtRequest struct {
	PersonID int `json:"id"`
}

// FillCourt method
func FillCourt(writer http.ResponseWriter, request *http.Request) {
	f := functionFillCourt
	ctx := request.Context()

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	str := mux.Vars(request)["id"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		message := fmt.Sprintf("the key [%s] is not an int", str)
		d := Dump(f, request, message)
		d.AddString("id", str)
		writeResponseMessage(writer, request, http.StatusBadRequest, message)
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

	positions, err := model.FillCourtTx(ctx, db, courtID)
	if err != nil {
		message := "problem filling court"
		d := Dump(f, request, message)
		d.AddString("courtID", fmt.Sprintf("%d", courtID))
		writeResponseError(writer, request, err)
		return
	}

	response := struct {
		Message   string           `json:"message"`
		Positions []model.Position `json:"positions"`
	}{
		Message:   "ok",
		Positions: positions,
	}
	writeResponseObject(writer, request, http.StatusOK, response)
}
