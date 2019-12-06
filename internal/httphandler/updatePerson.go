package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

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
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	f := functionUpdatePerson

	sess, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(w, r, http.StatusInternalServerError, "", "Error")
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request UpdatePersonRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(r)["id"]

	f.DebugVerbose("ID:     %s", id)
	f.DebugVerbose("Person: %v", request.Person)

	err = model.UpdatePerson(userID, id, request.Person)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
