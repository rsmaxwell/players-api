package httphandler

import (
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

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Person map[string]interface{} `json:"person"`
}

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(writer http.ResponseWriter, request *http.Request) {
	f := functionUpdatePerson
	ctx := request.Context()

	userID, err := checkAuthenticated(request)
	if err != nil {
		message := "Not Authenticated"
		Dump(f, request, message)
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

	limitedReader := &io.LimitedReader{R: request.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		message := "Problem reading body"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	DebugRequestBody(f, request, b)

	var updatePersonRequest UpdatePersonRequest
	err = json.Unmarshal(b, &updatePersonRequest)
	if err != nil {
		message := "Problem unmarshalling body"
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	str := mux.Vars(request)["id"]
	personID, err := strconv.Atoi(str)
	if err != nil {
		message := fmt.Sprintf("the key [%s] is not an int", str)
		DebugInfo(f, request, message)
		Dump(f, request, message)
		writeResponseMessage(writer, request, http.StatusBadRequest, message)
	}

	user := model.FullPerson{ID: userID}
	err = user.LoadPerson(ctx, db)
	if err != nil {
		message := fmt.Sprintf("Could not load person [%d]", userID)
		DebugVerbose(f, request, message)
		DumpError(f, request, err, message)
		writeResponseMessage(writer, request, http.StatusInternalServerError, message)
	}

	if userID == personID {
		err = user.CanEditSelf()
		if err != nil {
			message := "Forbidden: Not allowed to edit self"
			DebugVerbose(f, request, message)
			DumpError(f, request, err, message)
			writeResponseMessage(writer, request, http.StatusForbidden, "Forbidden")
		}
	} else {
		err = user.CanEditOtherPeople()
		if err != nil {
			message := "Forbidden: Not allowed to edit other people"
			DebugVerbose(f, request, message)
			DumpError(f, request, err, message)
			writeResponseMessage(writer, request, http.StatusForbidden, "Forbidden")
		}
	}

	err = model.UpdatePersonFieldsTx(db, personID, updatePersonRequest.Person)
	if err != nil {
		message := fmt.Sprintf("problem updating person fields: userID: %d", userID)
		d := DumpError(f, request, err, message)
		d.AddObject("request.Person", updatePersonRequest.Person)
		writeResponseError(writer, request, err)
	}

	writeResponseMessage(writer, request, http.StatusOK, "ok")
}
