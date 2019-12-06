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
func ListPeople(w http.ResponseWriter, r *http.Request) {
	f := functionListPeople

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

	var request ListPeopleRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		writeResponseMessage(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	f.DebugVerbose("Filter:    %s", request.Filter)

	listOfPeople, err := model.ListPeople(request.Filter)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", ListPeopleResponse{
		People: listOfPeople,
	})
}
