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

	session, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, err.Error())
		common.MetricsData.ServerError++
		return
	}
	defer session.SessionRelease(rw)
	value := session.Get("id")
	if value == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		return
	}

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r PostMoveRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	err = model.PostMove(&r.Source, &r.Target, r.Players)
	if err != nil {
		errorHandler(rw, req, err)
		common.MetricsData.ClientError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
