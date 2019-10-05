package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/session"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

// UpdateCourtRequest structure
type UpdateCourtRequest struct {
	Token string                 `json:"token"`
	Court map[string]interface{} `json:"court"`
}

// UpdateCourt method
func UpdateCourt(rw http.ResponseWriter, req *http.Request, id string) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r UpdateCourtRequest
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

	if !model.PersonCanUpdateCourt(session.UserID) {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	ref := common.Reference{Type: "court", ID: id}
	err = model.UpdateCourt(&ref, r.Court)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
