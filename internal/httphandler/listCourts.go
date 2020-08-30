package httphandler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// ListCourtsResponse structure
type ListCourtsResponse struct {
	Message string `json:"message"`
	Courts  []int  `json:"courts"`
}

var (
	functionListCourts = debug.NewFunction(pkg, "ListCourts")
)

// ListCourts method
func ListCourts(w http.ResponseWriter, r *http.Request) {
	f := functionListCourts
	f.DebugVerbose("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	object := r.Context().Value(ContextDatabaseKey)
	db, ok := object.(*sql.DB)
	if !ok {
		err = fmt.Errorf("Unexpected context type")
		writeResponseError(w, r, err)
		return
	}

	listOfCourts, err := model.ListCourts(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", ListCourtsResponse{
		Message: "ok",
		Courts:  listOfCourts,
	})
}
