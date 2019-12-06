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
func CreateCourt(w http.ResponseWriter, r *http.Request) {
	f := functionCreateCourt

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

	var request CreateCourtRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	id, err := model.CreateCourt(&request.Court)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseMessage(w, r, http.StatusOK, "", "ok")
	json.NewEncoder(w).Encode(CreateCourtResponse{
		ID: id,
	})
}
