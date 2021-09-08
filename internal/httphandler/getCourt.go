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
func GetCourt(writer http.ResponseWriter, request *http.Request) {
	f := functionGetCourt
	ctx := request.Context()

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

	object := request.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}

	var c model.Court
	c.ID = id
	err = c.LoadCourt(ctx, db)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseObject(writer, request, http.StatusOK, GetCourtResponse{
		Message: "ok",
		Court:   c,
	})
}
