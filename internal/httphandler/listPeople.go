package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
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

// ListPeople method
func ListPeople(rw http.ResponseWriter, req *http.Request) {

	log.Printf("ListPeople:")
	log.Printf("    Method: %s", req.Method)
	log.Printf("    Proto:  %s", req.Proto)
	log.Printf("    Host:   %s", req.Host)
	log.Printf("    URL:    %s", req.URL)

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	var r ListPeopleRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	log.Printf("    Filter:    %s", r.Filter)

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
