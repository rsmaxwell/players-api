package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdatePersonPlayer method
func UpdatePersonPlayer(token string, id string, player bool) error {
	f := debug.NewFunction(pkg, "UpdatePersonPlayer")
	f.DebugVerbose("")

	session := session.LookupToken(token)
	if session == nil {
		f.Verbosef("Not Authorised")
		return codeerror.NewUnauthorized("Not Authorised")
	}

	if !person.CanUpdatePersonPlayer(session.UserID, id) {
		f.Verbosef("Not Authorised")
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := person.UpdatePlayer(id, player)
	if err != nil {
		return err
	}

	return nil
}
