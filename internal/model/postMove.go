package model

import (
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/commands"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// PostMove method
func PostMove(token string, source, target *common.Reference, players []string) error {

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := commands.Move(source, target, players)
	if err != nil {
		return err
	}

	return nil
}
