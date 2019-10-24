package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionUpdatePersonPlayer = debug.NewFunction(pkg, "UpdatePersonPlayer")
)

// UpdatePersonPlayer method
func UpdatePersonPlayer(token string, id string, player bool) error {
	f := functionUpdatePersonPlayer
	f.DebugVerbose("token: %s, id: %s, player: %t", token, id, player)

	session := session.LookupToken(token)
	if session == nil {
		f.Verbosef("Not Authorised")
		return codeerror.NewUnauthorized("Not Authorised")
	}

	p, err := person.Load(session.UserID)
	if err != nil {
		return err
	}

	if !p.CanUpdatePersonPlayer() {
		f.Verbosef("Not Authorised")
		common.MetricsData.ClientError++
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err = person.UpdatePlayer(id, player)
	if err != nil {
		return err
	}

	return nil
}
