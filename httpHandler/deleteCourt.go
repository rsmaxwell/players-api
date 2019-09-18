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

// DeleteCourtRequest structure
type DeleteCourtRequest struct {
	Token string `json:"token"`
}

// DeleteCourt method
func DeleteCourt(rw http.ResponseWriter, req *http.Request, id string) {

	var r DeleteCourtRequest

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

	ok, err := court.Delete(id)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not delete court[%s]", id))
		clientError++
		return
	}

	setHeaders(rw, req)
	if ok {
		WriteResponse(rw, http.StatusOK, "ok")
	} else {
		WriteResponse(rw, http.StatusNotFound, "ok")
	}
}
