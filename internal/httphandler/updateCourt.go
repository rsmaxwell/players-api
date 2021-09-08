package httphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdateCourtRequest structure
type UpdateCourtRequest struct {
	Court map[string]interface{} `json:"court"`
}

var (
	functionUpdateCourt = debug.NewFunction(pkg, "UpdateCourt")
)

// UpdateCourt method
func UpdateCourt(writer http.ResponseWriter, request *http.Request) {
	f := functionUpdateCourt

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

	user := model.FullPerson{ID: userID}
	err = user.LoadPerson(context.Background(), db)
	if err != nil {
		message := fmt.Sprintf("Could not load person [%d]", userID)
		DumpError(f, request, err, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
		return
	}
	err = user.CanEditCourt()
	if err != nil {
		DebugVerbose(f, request, fmt.Sprintf("Person [%d] is not allowed to edit court", userID))
		writeResponseMessage(writer, request, http.StatusForbidden, "Forbidden")
	}

	limitedReader := &io.LimitedReader{R: request.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	DebugRequestBody(f, request, b)

	var updateCourtRequest UpdateCourtRequest
	err = json.Unmarshal(b, &updateCourtRequest)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	str := mux.Vars(request)["id"]
	courtID, err := strconv.Atoi(str)
	if err != nil {
		writeResponseMessage(writer, request, http.StatusBadRequest, fmt.Sprintf("the key [%s] is not an int", str))
		return
	}

	err = model.UpdateCourtFieldsTx(db, courtID, updateCourtRequest.Court)
	if err != nil {
		writeResponseError(writer, request, err)
		return
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
