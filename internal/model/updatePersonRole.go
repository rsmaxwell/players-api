package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdatePersonRole method
func UpdatePersonRole(token string, id string, role string) error {
	f := debug.NewFunction(pkg, "UpdatePersonRole")
	f.DebugVerbose("")

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	if !person.CanUpdatePersonRole(session.UserID, id) {
		f.Verbosef("Unauthorized")
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := person.UpdateRole(id, role)
	if err != nil {
		return err
	}

	return nil
}
