package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
)

// DeletePerson method
func DeletePerson(token, id string) error {
	f := functionDeletePerson
	f.DebugVerbose("token: %s id:%s", token, id)

	session := session.LookupToken(token)
	if session == nil {
		f.Dump("token not found")
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := person.Remove(id)
	if err != nil {
		f.Dump("could not remove person: %v", err)
		return err
	}

	return nil
}
