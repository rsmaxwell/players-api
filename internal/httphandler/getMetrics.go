package httphandler

import (
	"database/sql"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetMetricsResponse structure
type GetMetricsResponse struct {
	Message string        `json:"message"`
	Data    model.Metrics `json:"metrics"`
}

var (
	functionGetMetrics = debug.NewFunction(pkg, "GetMetrics")
)

// GetMetrics method
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	f := functionGetMetrics
	f.DebugAPI("")

	userID, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		message := "unexpected context type"
		f.Dump(message)
		writeResponseMessage(w, r, http.StatusInternalServerError, message)
		return
	}

	p := model.FullPerson{ID: userID}
	err = p.LoadPerson(db)
	if err != nil {
		f.Dump("Could not load the logged on user: %d", userID)
		writeResponseMessage(w, r, http.StatusInternalServerError, "Not Authorized")
		return
	}

	err = p.CanGetMetrics()
	if err != nil {
		f.DebugVerbose("unauthorized person[%d] attempted to get metrics", userID)
		writeResponseMessage(w, r, http.StatusForbidden, "Forbidden")
		return
	}

	writeResponseObject(w, r, http.StatusOK, GetMetricsResponse{
		Message: "ok",
		Data:    model.MetricsData,
	})
}
