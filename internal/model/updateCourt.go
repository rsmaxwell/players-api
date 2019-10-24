package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionUpdateCourt = debug.NewFunction(pkg, "UpdateCourt")
)

// UpdateCourt method
func UpdateCourt(token string, id string, fields map[string]interface{}) error {
	f := functionUpdateCourt
	f.DebugVerbose("token: %s, id: %s, fields: %v", token, id, fields)

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	p, err := person.Load(session.UserID)
	if err != nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	if !p.CanUpdateCourt() {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	ref := &common.Reference{Type: "court", ID: id}
	err = court.Update(ref, fields)
	if err != nil {
		return err
	}

	return nil
}
