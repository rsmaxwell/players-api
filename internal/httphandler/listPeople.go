package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// ListPeopleRequest structure
type ListPeopleRequest struct {
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

	var r ListPeopleRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		writeResponseMessage(rw, req, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugVerbose("Filter:    %s", r.Filter)

	listOfPeople, err := model.ListPeople(r.Filter)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", ListPeopleResponse{
		People: listOfPeople,
	})
}
