package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionDeleteCourt = debug.NewFunction(pkg, "DeleteCourt")
)

// DeleteCourt method
func DeleteCourt(token string, id string) error {
	f := functionDeleteCourt
	f.DebugVerbose("token: %s, id:%s", token, id)

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := court.Remove(id)
	if err != nil {
		return err
	}

	return nil
}
