package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/person"

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
func UpdateCourt(w http.ResponseWriter, r *http.Request) {
	f := functionUpdateCourt

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

	p, err := person.Load(userID)
	if err != nil {
		f.Dump("Could not load the logged on user[%s]: %v", userID, err)
		writeResponseMessage(w, r, http.StatusInternalServerError, "", err.Error())
		return
	}
	if !p.CanUpdateCourt() {
		writeResponseMessage(w, r, http.StatusForbidden, "", "Forbidden")
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request UpdateCourtRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdateCourt(id, request.Court)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
