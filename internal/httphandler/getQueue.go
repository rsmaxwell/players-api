package httphandler

import (
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetQueueResponse structure
type GetQueueResponse struct {
	Queue queue.Queue `json:"queue"`
}

// GetQueue method
func GetQueue(rw http.ResponseWriter, req *http.Request) {

	_, err := checkAccessToken(req)
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	q, err := model.GetQueue()
	if err != nil {
		writeResponseError(rw, req, err)
		return
	}

	writeResponseObject(rw, req, http.StatusOK, "", GetQueueResponse{
		Queue: *q,
	})
}
