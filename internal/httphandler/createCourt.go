package httphandler

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// CreateCourtRequest structure
type CreateCourtRequest struct {
	Court model.Court `json:"court"`
}

// CreateCourtResponse structure
type CreateCourtResponse struct {
	Message string      `json:"message"`
	Court   model.Court `json:"court"`
}

var (
	functionCreateCourt = debug.NewFunction(pkg, "CreateCourt")
)

// CreateCourt method
func CreateCourt(w http.ResponseWriter, r *http.Request) {
	f := functionCreateCourt
	f.DebugAPI("")

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request CreateCourtRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
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

	c := model.Court{Name: request.Court.Name}
	err = c.SaveCourt(db)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(CreateCourtResponse{
		Message: "ok",
		Court:   c,
	})
}
