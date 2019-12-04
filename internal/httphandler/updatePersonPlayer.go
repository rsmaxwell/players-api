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
func UpdatePersonPlayer(rw http.ResponseWriter, req *http.Request) {
	f := functionUpdatePersonPlayer

	sess, err := checkAuthenticated(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	userID, ok := sess.Values["userID"].(string)
	if !ok {
		f.DebugVerbose("could not get 'userID' from the session")
		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "Error")
		return
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var r UpdatePersonPlayerRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	id := mux.Vars(req)["id"]
	f.DebugVerbose("ID: %s", id)

	err = model.UpdatePersonPlayer(userID, id, r.Player)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
