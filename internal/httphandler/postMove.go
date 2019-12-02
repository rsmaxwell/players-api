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
func PostMove(rw http.ResponseWriter, req *http.Request) {
	f := functionPostMove

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugRequestBody(b)

	var r PostMoveRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	err = model.PostMove(&r.Source, &r.Target, r.Players)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
}
