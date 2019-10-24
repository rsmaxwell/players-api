package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(token string, id string, fields map[string]interface{}) error {
	f := functionUpdatePerson
	f.DebugVerbose("token: %s, id: %s, fields: %v", token, id, fields)

	session := session.LookupToken(token)
	if session == nil {
		f.Verbosef("Not Authorised")
		return codeerror.NewUnauthorized("Not Authorised")
	}

	p, err := person.Load(session.UserID)
	if err != nil {
		return err
	}

	if !p.CanUpdatePerson(session.UserID, id) {
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err = person.Update(id, fields)
	if err != nil {
		return err
	}

	return nil
}
