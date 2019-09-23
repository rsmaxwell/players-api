package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/queue"
	"github.com/rsmaxwell/players-api/session"
)

// GetQueueRequest structure
type GetQueueRequest struct {
	Token string `json:"token"`
}

// GetQueueResponse structure
type GetQueueResponse struct {
	Queue queue.Queue `json:"queue"`
}

// GetQueue method
func GetQueue(rw http.ResponseWriter, req *http.Request) {

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	var r GetQueueRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		clientError++
		return
	}

	session := session.LookupToken(r.Token)
	if session == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
		clientError++
		return
	}

	queue, err := queue.Load()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetQueueResponse{
		Queue: *queue,
	})
}
