package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
)

// GetCourtRequest structure
type GetCourtRequest struct {
	Token string `json:"token"`
}

// GetCourtResponse structure
type GetCourtResponse struct {
	Court model.Court `json:"court"`
}

// GetCourt method
func GetCourt(rw http.ResponseWriter, req *http.Request, id string) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r GetCourtRequest
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

	ref := &common.Reference{Type: "court", ID: id}
	court, err := model.LoadCourt(ref)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetCourtResponse{
		Court: *court,
	})
}
