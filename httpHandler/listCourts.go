package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/court"
	"github.com/rsmaxwell/players-api/session"
)

// ListCourtsRequest structure
type ListCourtsRequest struct {
	Token string `json:"token"`
}

// List courts Response
type listCourtsResponse struct {
	Courts []string `json:"courts"`
}

// ListCourts method
func ListCourts(rw http.ResponseWriter, req *http.Request) {

	var r ListCourtsRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
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

	listOfCourts, err := court.List()
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, "Error getting list of courts")
		serverError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listCourtsResponse{
		Courts: listOfCourts,
	})
}
