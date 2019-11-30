package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionUpdatePersonRole = debug.NewFunction(pkg, "UpdatePersonRole")
)

// UpdatePersonRole method
func UpdatePersonRole(userID string, id string, role string) error {
	f := functionUpdatePersonRole
	f.DebugVerbose("userID: %s, id: %s, role: %s", userID, id, role)

	p, err := person.Load(userID)
	if err != nil {
		return err
	}

	if !p.CanUpdatePersonRole(userID, id) {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	err = person.UpdateRole(id, role)
	if err != nil {
		return err
	}

	return nil
}
