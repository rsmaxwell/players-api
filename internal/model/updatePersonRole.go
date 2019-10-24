package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionUpdatePersonRole = debug.NewFunction(pkg, "UpdatePersonRole")
)

// UpdatePersonRole method
func UpdatePersonRole(token string, id string, role string) error {
	f := functionUpdatePersonRole
	f.DebugVerbose("token: %s, id: %s, role: %s", token, id, role)

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	p, err := person.Load(session.UserID)
	if err != nil {
		return err
	}

	if !p.CanUpdatePersonRole(session.UserID, id) {
		f.Verbosef("Unauthorized")
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err = person.UpdateRole(id, role)
	if err != nil {
		return err
	}

	return nil
}
