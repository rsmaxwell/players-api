package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/session"
)

// GetMetrics method
func GetMetrics(token string) error {

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !person.CanGetMetrics(session.UserID) {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	return nil
}
