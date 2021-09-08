package httphandler

import (
	"database/sql"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionListCourts = debug.NewFunction(pkg, "ListCourts")
)

// ListCourts method
func ListCourts(writer http.ResponseWriter, request *http.Request) {
	f := functionListCourts

	_, err := checkAuthenticated(request)
	if err != nil {
		message := "Authorisation problem"
		Dump(f, request, message)
		writeResponseError(writer, request, err)
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

	listOfCourts, err := model.ListCourtsTx(db)
	if err != nil {
		message := "Problem listing courts"
		Dump(f, request, message)
		writeResponseError(writer, request, err)
		return
	}

	writeResponseObject(writer, request, http.StatusOK, listOfCourts)
}
