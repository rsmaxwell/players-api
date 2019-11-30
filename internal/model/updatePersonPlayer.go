package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionUpdatePersonPlayer = debug.NewFunction(pkg, "UpdatePersonPlayer")
)

// UpdatePersonPlayer method
func UpdatePersonPlayer(userID string, id string, player bool) error {
	f := functionUpdatePersonPlayer
	f.DebugVerbose("userID: %s, id: %s, player: %t", userID, id, player)

	p, err := person.Load(userID)
	if err != nil {
		return err
	}

	if !p.CanUpdatePersonPlayer() {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	err = person.UpdatePlayer(id, player)
	if err != nil {
		return err
	}

	return nil
}
