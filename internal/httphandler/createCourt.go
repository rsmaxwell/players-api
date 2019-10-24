package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/model"
)

// CreateCourtRequest structure
type CreateCourtRequest struct {
	Token string      `json:"token"`
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

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r CreateCourtRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	id, err := model.CreateCourt(r.Token, &r.Court)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(CreateCourtResponse{
		ID: id,
	})
}
