package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdatePerson method
func UpdatePerson(token string, id string, fields map[string]interface{}) error {
	f := debug.NewFunction(pkg, "UpdatePerson")
	f.DebugVerbose("")

	session := session.LookupToken(token)
	if session == nil {
		f.Verbosef("Not Authorised")
		return codeerror.NewUnauthorized("Not Authorised")
	}

	if !person.CanUpdatePerson(session.UserID, id) {
		f.Verbosef("Not Authorised")
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := person.Update(id, fields)
	if err != nil {
		return err
	}

	return nil
}
