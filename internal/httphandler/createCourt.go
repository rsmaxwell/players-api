package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/model"
)

// CreateCourtRequest structure
type CreateCourtRequest struct {
	Court court.Court `json:"court"`
}

// CreateCourtResponse structure
type CreateCourtResponse struct {
	ID string `json:"id"`
}

var (
	functionCreateCourt = debug.NewFunction(pkg, "CreateCourt")
)

// CreateCourt method
func CreateCourt(rw http.ResponseWriter, req *http.Request) {
	f := functionCreateCourt

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

	var r CreateCourtRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	id, err := model.CreateCourt(&r.Court)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseMessage(rw, req, http.StatusOK, "", "ok")
	json.NewEncoder(rw).Encode(CreateCourtResponse{
		ID: id,
	})
}
