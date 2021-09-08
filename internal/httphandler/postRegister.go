package httphandler

import (
	"context"
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
func Register(writer http.ResponseWriter, request *http.Request) {
	f := functionRegister

	limitedReader := &io.LimitedReader{R: request.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	DebugRequestBody(f, request, b)

	var registrationRequest model.Registration
	err = json.Unmarshal(b, &registrationRequest)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	// f.DebugVerbose("FirstName: %s", registrationRequest.FirstName)
	// f.DebugVerbose("LastName:  %s", registrationRequest.LastName)
	// f.DebugVerbose("Knownas:   %s", registrationRequest.Knownas)
	// f.DebugVerbose("Email:     %s", registrationRequest.Email)
	// f.DebugVerbose("Phone:     %s", registrationRequest.Phone)
	// f.DebugVerbose("Password:  %s", registrationRequest.Password)

	p, err := registrationRequest.ToPerson()
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

	err = p.SavePerson(context.Background(), db)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				err = codeerror.NewBadRequest(pgerr.Message)
			}
		}

		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
