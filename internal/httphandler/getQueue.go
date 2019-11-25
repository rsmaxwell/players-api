package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

// GetQueueResponse structure
type GetQueueResponse struct {
	Queue queue.Queue `json:"queue"`
}

// GetQueue method
func GetQueue(rw http.ResponseWriter, req *http.Request) {

	session, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, err.Error())
		common.MetricsData.ServerError++
		return
	}
	defer session.SessionRelease(rw)
	userID := session.Get("id")
	if userID == nil {
		WriteResponse(rw, http.StatusUnauthorized, "Not Authorized")
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
