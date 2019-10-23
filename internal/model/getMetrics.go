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

	p, err := person.Load(session.UserID)
	if err != nil {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CanGetMetrics() {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	return nil
}
