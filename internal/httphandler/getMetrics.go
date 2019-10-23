package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetMetricsRequest structure
type GetMetricsRequest struct {
	Token string `json:"token"`
}

// GetMetricsResponse structure
type GetMetricsResponse struct {
	Data common.Metrics `json:"data"`
}

// GetMetrics method
func GetMetrics(rw http.ResponseWriter, req *http.Request) {
	f := debug.NewFunction(pkg, "GetMetrics")

	limitedReader := io.LimitReader(req.Body, 100*1024)
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r GetMetricsRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	err = model.GetMetrics(r.Token)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetMetricsResponse{
		Data: common.MetricsData,
	})
}
