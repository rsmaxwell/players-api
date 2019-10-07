package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// GetQueue method
func GetQueue(token string) (*queue.Queue, error) {

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	ref := common.Reference{Type: "queue", ID: ""}
	q, err := queue.Load(&ref)
	if err != nil {
		return nil, err
	}

	return q, nil
}
