package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionUpdatePerson = debug.NewFunction(pkg, "UpdatePerson")
)

// UpdatePerson method
func UpdatePerson(userID string, id string, fields map[string]interface{}) error {
	f := functionUpdatePerson
	f.DebugVerbose("userId: %s, id: %s, fields: %v", userID, id, fields)

	p, err := person.Load(userID)
	if err != nil {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CanUpdatePerson(userID, id) {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	err = person.Update(id, fields)
	if err != nil {
		return err
	}

	return nil
}
