package httphandler

import (
	"context"
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

var (
	functionCreateCourt = debug.NewFunction(pkg, "CreateCourt")
)

// CreateCourt method
func CreateCourt(writer http.ResponseWriter, request *http.Request) {
	f := functionCreateCourt

	_, err := checkAuthenticated(request)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	limitedReader := &io.LimitedReader{R: request.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	DebugRequestBody(f, request, b)

	var createCourtRequest CreateCourtRequest
	err = json.Unmarshal(b, &createCourtRequest)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
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

	c := model.Court{Name: createCourtRequest.Court.Name}
	err = c.SaveCourt(context.Background(), db)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseObject(writer, request, http.StatusOK, c)
}
