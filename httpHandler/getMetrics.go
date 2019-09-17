package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/session"
)

// GetMetricsRequest structure
type GetMetricsRequest struct {
	Token string `json:"token"`
}

// metrics Response
type metricsResponse struct {
	ClientSuccess             int `json:"clientSuccess"`
	ClientError               int `json:"clientError"`
	ClientAuthenticationError int `json:"clientAuthenticationError"`
	ServerError               int `json:"serverError"`
}

// GetMetrics method
func GetMetrics(rw http.ResponseWriter, req *http.Request) {

	var r GetMetricsRequest

	// limitedReader := &io.LimitedReader{R: req.Body, N: 100 * 1024}

	limitedReader := io.LimitReader(req.Body, 100*1024)
	b, err := ioutil.ReadAll(limitedReader)

	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse data"))
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(metricsResponse{
		ClientSuccess:             clientSuccess,
		ClientError:               clientError,
		ClientAuthenticationError: clientAuthenticationError,
		ServerError:               serverError,
	})
}
