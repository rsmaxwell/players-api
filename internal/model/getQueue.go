package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionGetQueue = debug.NewFunction(pkg, "GetQueue")
)

// GetQueue method
func GetQueue() (*queue.Queue, error) {
	f := functionGetQueue
	f.DebugVerbose("")

	ref := common.Reference{Type: "queue", ID: ""}
	q, err := queue.Load(&ref)
	if err != nil {
		return nil, err
	}

	return q, nil
}
