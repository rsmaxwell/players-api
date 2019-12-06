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

// UpdatePersonRoleRequest structure
type UpdatePersonRoleRequest struct {
	Role string `json:"role"`
}

var (
	functionUpdatePersonRole = debug.NewFunction(pkg, "UpdatePersonRole")
)

// UpdatePersonRole method
func UpdatePersonRole(w http.ResponseWriter, r *http.Request) {
	f := functionUpdatePersonRole

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

	var request UpdatePersonRoleRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdatePersonRole(userID, id, request.Role)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
