package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// PostMoveRequest structure
type PostMoveRequest struct {
	Source  common.Reference `json:"source"`
	Target  common.Reference `json:"target"`
	Players []string         `json:"players"`
}

var (
	functionPostMove = debug.NewFunction(pkg, "PostMove")
)

// PostMove method
func PostMove(w http.ResponseWriter, r *http.Request) {
	f := functionPostMove

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var request PostMoveRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	err = model.PostMove(&request.Source, &request.Target, request.Players)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
}
