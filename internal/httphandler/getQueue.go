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
func GetQueue(w http.ResponseWriter, r *http.Request) {

	_, err := checkAuthenticated(r)
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	q, err := model.GetQueue()
	if err != nil {
		writeResponseError(w, r, err)
		return
	}

	writeResponseObject(w, r, http.StatusOK, "", GetQueueResponse{
		Queue: *q,
	})
}
