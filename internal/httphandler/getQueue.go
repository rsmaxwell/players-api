package httphandler

import (
	"encoding/json"
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

	_, err := checkAuthToken(req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	q, err := model.GetQueue()
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(GetQueueResponse{
		Queue: *q,
	})
}
