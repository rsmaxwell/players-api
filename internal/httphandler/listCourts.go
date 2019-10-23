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

// ListCourtsRequest structure
type ListCourtsRequest struct {
	Token string `json:"token"`
}

// ListCourtsResponse structure
type ListCourtsResponse struct {
	Courts []string `json:"courts"`
}

// ListCourts method
func ListCourts(rw http.ResponseWriter, req *http.Request) {
	f := debug.NewFunction(pkg, "ListCourts")

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r ListCourtsRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	listOfCourts, err := model.ListCourts(r.Token)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(ListCourtsResponse{
		Courts: listOfCourts,
	})
}
