package httphandler

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/rsmaxwell/players-api/internal/debug"
)

// // RegisterRequest structure
// type RegisterRequest struct {
// 	Registration model.Registration `json:"register"`
// }

var (
	functionRegister = debug.NewFunction(pkg, "Register")
)

// Register method
func Register(w http.ResponseWriter, r *http.Request) {
	f := functionRegister
	f.DebugAPI("")

	if r.Method == http.MethodOptions {
		writeResponseMessage(w, r, http.StatusOK, "ok")
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request model.Registration
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// f.DebugVerbose("FirstName: %s", request.FirstName)
	// f.DebugVerbose("LastName:  %s", request.LastName)
	// f.DebugVerbose("Knownas:   %s", request.Knownas)
	// f.DebugVerbose("Email:     %s", request.Email)
	// f.DebugVerbose("Phone:     %s", request.Phone)
	// f.DebugVerbose("Password:  %s", request.Password)

	p, err := request.ToPerson()
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

	err = p.SavePerson(db)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				err = codeerror.NewBadRequest(pgerr.Message)
			}
		}

		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "ok")
}
