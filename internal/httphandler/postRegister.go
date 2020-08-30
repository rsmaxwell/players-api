package httphandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/rsmaxwell/players-api/internal/debug"
)

// RegisterRequest structure
type RegisterRequest struct {
	Registration model.Registration `json:"registration"`
}

var (
	functionRegister = debug.NewFunction(pkg, "Register")
)

// Register method
func Register(w http.ResponseWriter, r *http.Request) {
	f := functionRegister

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request RegisterRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugVerbose("FirstName: %s", request.Registration.FirstName)
	f.DebugVerbose("LastName:  %s", request.Registration.LastName)
	f.DebugVerbose("DisplayName: %s", request.Registration.DisplayName)
	f.DebugVerbose("UserName:  %s", request.Registration.UserName)
	f.DebugVerbose("Email:     %s", request.Registration.Email)
	f.DebugVerbose("Phone:     %s", request.Registration.Phone)
	f.DebugVerbose("Password:  %s", request.Registration.Password)

	p, err := request.Registration.ToPerson()
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

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
