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

// UpdatePersonPlayerRequest structure
type UpdatePersonPlayerRequest struct {
	Token  string `json:"token"`
	Player bool   `json:"player"`
}

var (
	functionUpdatePersonPlayer = debug.NewFunction(pkg, "UpdatePersonPlayer")
)

// UpdatePersonPlayer method
func UpdatePersonPlayer(w http.ResponseWriter, r *http.Request) {
	f := functionUpdatePersonPlayer

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

	var request UpdatePersonPlayerRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(r)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdatePersonPlayer(userID, id, request.Player)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
