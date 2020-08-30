package httphandler

import (
	"database/sql"
	"fmt"
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

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	userID, ok := sess.Values["userID"].(int)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("Unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	p := model.Person{ID: userID}
	err = p.LoadPerson(db)
	if err != nil {
		f.Dump("Could not load the logged on user: %d", userID)
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Not Authorized")
		return
	}

	err = p.CanGetMetrics()
	if err != nil {
		f.DebugVerbose("unauthorized person[%d] attempted to get metrics", userID)
		writeResponseMessage(w, r, http.StatusForbidden, "", "Forbidden")
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetMetricsResponse{
		Message: "ok",
		Data:    model.MetricsData,
	})
}
