package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionGetCourt = debug.NewFunction(pkg, "GetCourt")
)

// GetCourt method
func GetCourt(token, id string) (*court.Court, error) {
	f := functionGetCourt
	f.DebugVerbose("token: %s id:%s", token, id)

	session := session.LookupToken(token)
	if session == nil {
		f.Dump("token not found")
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	ref := &common.Reference{Type: "court", ID: id}
	court, err := court.Load(ref)
	if err != nil {
		f.Dump("could not load court: %v", err)
		return nil, err
	}

	return court, nil
}
