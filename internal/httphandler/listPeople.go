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

// ListPeopleRequest structure
type ListPeopleRequest struct {
	Token  string   `json:"token"`
	Filter []string `json:"filter"`
}

// ListPeopleResponse structure
type ListPeopleResponse struct {
	People []string `json:"people"`
}

var (
	functionListPeople = debug.NewFunction(pkg, "ListPeople")
)

// ListPeople method
func ListPeople(rw http.ResponseWriter, req *http.Request) {
	f := functionListPeople

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugRequestBody(b)

	var r ListPeopleRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	f.DebugVerbose("Filter:    %s", r.Filter)

	listOfPeople, err := model.ListPeople(r.Token, r.Filter)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(ListPeopleResponse{
		People: listOfPeople,
	})
}
