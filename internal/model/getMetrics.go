package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionGetMetrics = debug.NewFunction(pkg, "GetMetrics")
)

// GetMetrics method
func GetMetrics(token string) error {
	f := functionGetMetrics
	f.DebugVerbose("token: %s", token)

	session := session.LookupToken(token)
	if session == nil {
		f.Dump("token not found")
		return codeerror.NewUnauthorized("Not Authorized")
	}

	p, err := person.Load(session.UserID)
	if err != nil {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CanGetMetrics() {
		f.Dump("person not authorized to get metrics")
		return codeerror.NewUnauthorized("Not Authorized")
	}

	return nil
}
