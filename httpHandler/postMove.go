package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/commands"
	"github.com/rsmaxwell/players-api/destination"
	"github.com/rsmaxwell/players-api/session"
)

// PostMoveRequest structure
type PostMoveRequest struct {
	Token   string                `json:"token"`
	Source  destination.Reference `json:"source"`
	Target  destination.Reference `json:"target"`
	Players []string              `json:"players"`
}

// PostMove method
func PostMove(rw http.ResponseWriter, req *http.Request) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r PostMoveRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	err = commands.Move(&r.Source, &r.Target, r.Players)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
