package httphandler

import (
	"database/sql"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

var (
	functionGetMetrics = debug.NewFunction(pkg, "GetMetrics")
)

// GetMetrics method
func GetMetrics(writer http.ResponseWriter, request *http.Request) {
	f := functionGetMetrics
	ctx := request.Context()

	userID, err := checkAuthenticated(request)
	if err != nil {
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

	p := model.FullPerson{ID: userID}
	err = p.LoadPerson(ctx, db)
	if err != nil {
		Dump(f, request, "Could not load the logged on user: %d", userID)
		writeResponseMessage(writer, request, http.StatusInternalServerError, "Not Authorized")
		return
	}

	err = p.CanGetMetrics()
	if err != nil {
		DebugVerbose(f, request, "unauthorized person[%d] attempted to get metrics", userID)
		writeResponseMessage(writer, request, http.StatusForbidden, "Forbidden")
		return
	}

	writeResponseObject(writer, request, http.StatusOK, model.MetricsData)
}
